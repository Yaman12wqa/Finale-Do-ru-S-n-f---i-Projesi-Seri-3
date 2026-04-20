package cve

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/secscan/secscan/backend/internal/domain"
	"github.com/secscan/secscan/backend/internal/scanner"
	"github.com/secscan/secscan/backend/internal/utils"
)

const moduleName = "cve"

type Scanner struct {
	adapter AdvisoryAdapter
}

type Technology struct {
	Name     string `json:"name"`
	Version  string `json:"version,omitempty"`
	Evidence string `json:"evidence"`
}

type Advisory struct {
	ID         string `json:"id"`
	Technology string `json:"technology"`
	Title      string `json:"title"`
	Severity   string `json:"severity"`
	Affected   string `json:"affected,omitempty"`
	URL        string `json:"url"`
}

type AdvisoryAdapter interface {
	Find(technology Technology) []Advisory
}

type LocalAdapter struct {
	advisories map[string][]Advisory
}

func New() *Scanner {
	return &Scanner{adapter: NewLocalAdapter()}
}

func (s *Scanner) Name() string {
	return moduleName
}

func (s *Scanner) Scan(ctx context.Context, target utils.ValidatedTarget) domain.ModuleResult {
	result := scanner.NewResult(moduleName)
	client := utils.NewSafeHTTPClient(12 * time.Second)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target.CanonicalURL, nil)
	if err != nil {
		return scanner.Fail(moduleName, result.StartedAt, err)
	}
	req.Header.Set("User-Agent", "SecScan/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return scanner.Fail(moduleName, result.StartedAt, err)
	}
	body, err := utils.ReadLimited(resp, 512*1024)
	if err != nil {
		return scanner.Fail(moduleName, result.StartedAt, err)
	}

	technologies := detectTechnologies(resp.Header, body)
	allAdvisories := make([]Advisory, 0)
	score := 100

	for _, tech := range technologies {
		advisories := s.adapter.Find(tech)
		allAdvisories = append(allAdvisories, advisories...)
		if len(advisories) == 0 {
			continue
		}

		exact := hasExactAffectedVersion(tech, advisories)
		severity := "low"
		if exact {
			severity = "high"
			score -= 25
		} else {
			score -= 3
		}

		result.Findings = append(result.Findings, domain.Finding{
			ID:             "cve-" + strings.ToLower(tech.Name),
			Title:          "Technology has mapped advisory references",
			Description:    "A detected technology maps to local advisory data. This module is a pluggable mock adapter and does not confirm that the target is vulnerable.",
			Severity:       severity,
			Evidence:       fmt.Sprintf("%s %s from %s", tech.Name, tech.Version, tech.Evidence),
			Recommendation: "Verify the exact version and compare it with NVD, OSV, vendor advisories, or a future live adapter.",
		})
	}

	if len(technologies) == 0 {
		result.Findings = append(result.Findings, domain.Finding{
			ID:             "cve-no-tech-detected",
			Title:          "No common technologies detected",
			Description:    "The response did not expose recognizable technology headers or body markers from the local detector.",
			Severity:       "info",
			Recommendation: "Authenticated review and software bill of materials data provide better coverage than passive detection.",
		})
	}

	result.Metadata["detected_technologies"] = technologies
	result.Metadata["advisory_source"] = "local mock adapter"
	result.Metadata["advisories"] = allAdvisories

	scanner.Complete(&result, score)
	return result
}

func NewLocalAdapter() LocalAdapter {
	return LocalAdapter{advisories: map[string][]Advisory{
		"wordpress": {
			{
				ID:         "CVE-2022-21661",
				Technology: "WordPress",
				Title:      "WordPress core SQL injection advisory reference",
				Severity:   "high",
				Affected:   "<=5.8.2",
				URL:        "https://nvd.nist.gov/vuln/detail/CVE-2022-21661",
			},
		},
		"apache": {
			{
				ID:         "CVE-2021-41773",
				Technology: "Apache HTTP Server",
				Title:      "Path traversal and file disclosure in Apache 2.4.49",
				Severity:   "critical",
				Affected:   "2.4.49",
				URL:        "https://nvd.nist.gov/vuln/detail/CVE-2021-41773",
			},
		},
		"nginx": {
			{
				ID:         "CVE-2021-23017",
				Technology: "nginx",
				Title:      "Resolver vulnerability advisory reference",
				Severity:   "high",
				Affected:   "version-specific",
				URL:        "https://nvd.nist.gov/vuln/detail/CVE-2021-23017",
			},
		},
		"php": {
			{
				ID:         "CVE-2019-11043",
				Technology: "PHP-FPM",
				Title:      "PHP-FPM remote code execution advisory reference",
				Severity:   "critical",
				Affected:   "configuration-specific",
				URL:        "https://nvd.nist.gov/vuln/detail/CVE-2019-11043",
			},
		},
		"jquery": {
			{
				ID:         "CVE-2020-11023",
				Technology: "jQuery",
				Title:      "jQuery HTML parsing XSS advisory reference",
				Severity:   "medium",
				Affected:   "<3.5.0",
				URL:        "https://nvd.nist.gov/vuln/detail/CVE-2020-11023",
			},
		},
	}}
}

func (a LocalAdapter) Find(technology Technology) []Advisory {
	return append([]Advisory(nil), a.advisories[strings.ToLower(technology.Name)]...)
}

func detectTechnologies(headers http.Header, body string) []Technology {
	seen := map[string]Technology{}

	server := headers.Get("Server")
	if server != "" {
		detectHeaderTechnology(seen, "nginx", server, `(?i)nginx/?([0-9.]+)?`, "Server header")
		detectHeaderTechnology(seen, "apache", server, `(?i)apache/?([0-9.]+)?`, "Server header")
		detectHeaderTechnology(seen, "iis", server, `(?i)microsoft-iis/?([0-9.]+)?`, "Server header")
	}

	poweredBy := headers.Get("X-Powered-By")
	if poweredBy != "" {
		detectHeaderTechnology(seen, "php", poweredBy, `(?i)php/?([0-9.]+)?`, "X-Powered-By header")
		detectHeaderTechnology(seen, "express", poweredBy, `(?i)express`, "X-Powered-By header")
	}

	lowerBody := strings.ToLower(body)
	if strings.Contains(lowerBody, "wp-content") || strings.Contains(lowerBody, "wp-includes") {
		version := regexGroup(body, `(?i)generator" content="WordPress ([0-9.]+)`)
		seen["wordpress"] = Technology{Name: "wordpress", Version: version, Evidence: "WordPress body markers"}
	}
	if strings.Contains(lowerBody, "/_next/static/") {
		seen["next.js"] = Technology{Name: "next.js", Evidence: "Next.js static asset marker"}
	}
	if strings.Contains(lowerBody, "csrfmiddlewaretoken") {
		seen["django"] = Technology{Name: "django", Evidence: "Django CSRF token marker"}
	}
	if version := regexGroup(body, `(?i)jquery[.-]([0-9]+\.[0-9]+\.[0-9]+)`); version != "" {
		seen["jquery"] = Technology{Name: "jquery", Version: version, Evidence: "jQuery script reference"}
	}

	technologies := make([]Technology, 0, len(seen))
	for _, technology := range seen {
		technologies = append(technologies, technology)
	}
	return technologies
}

func detectHeaderTechnology(seen map[string]Technology, name, value, pattern, evidence string) {
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(value)
	if len(match) == 0 {
		return
	}
	version := ""
	if len(match) > 1 {
		version = match[1]
	}
	seen[name] = Technology{Name: name, Version: version, Evidence: evidence + ": " + value}
}

func regexGroup(value, pattern string) string {
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(value)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func hasExactAffectedVersion(technology Technology, advisories []Advisory) bool {
	if technology.Version == "" {
		return false
	}
	for _, advisory := range advisories {
		if advisory.Affected == technology.Version {
			return true
		}
		if strings.HasPrefix(advisory.Affected, "<=") && versionLessOrEqual(technology.Version, strings.TrimPrefix(advisory.Affected, "<=")) {
			return true
		}
	}
	return false
}

func versionLessOrEqual(current, maximum string) bool {
	currentParts := versionParts(current)
	maximumParts := versionParts(maximum)
	length := len(currentParts)
	if len(maximumParts) > length {
		length = len(maximumParts)
	}
	for i := 0; i < length; i++ {
		currentValue := 0
		maximumValue := 0
		if i < len(currentParts) {
			currentValue = currentParts[i]
		}
		if i < len(maximumParts) {
			maximumValue = maximumParts[i]
		}
		if currentValue < maximumValue {
			return true
		}
		if currentValue > maximumValue {
			return false
		}
	}
	return true
}

func versionParts(value string) []int {
	fields := strings.Split(value, ".")
	parts := make([]int, 0, len(fields))
	for _, field := range fields {
		parsed, err := strconv.Atoi(field)
		if err != nil {
			break
		}
		parts = append(parts, parsed)
	}
	return parts
}
