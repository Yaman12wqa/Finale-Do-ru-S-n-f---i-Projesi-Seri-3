package headers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/secscan/secscan/backend/internal/domain"
	"github.com/secscan/secscan/backend/internal/scanner"
	"github.com/secscan/secscan/backend/internal/utils"
)

const moduleName = "headers"

type Scanner struct{}

type headerRule struct {
	Name           string
	Recommendation string
}

var rules = []headerRule{
	{"Content-Security-Policy", "Define a strict Content-Security-Policy to reduce script injection impact."},
	{"Strict-Transport-Security", "Serve HTTPS and enable HSTS for production domains."},
	{"X-Frame-Options", "Set DENY or SAMEORIGIN unless framing is intentionally required."},
	{"X-Content-Type-Options", "Set X-Content-Type-Options to nosniff."},
	{"Referrer-Policy", "Use a privacy-preserving Referrer-Policy such as strict-origin-when-cross-origin."},
	{"Permissions-Policy", "Limit browser features with a Permissions-Policy header."},
	{"Cross-Origin-Opener-Policy", "Set COOP to isolate browsing contexts where possible."},
	{"Cross-Origin-Resource-Policy", "Set CORP to restrict cross-origin resource loading where appropriate."},
}

func New() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Name() string {
	return moduleName
}

func (s *Scanner) Scan(ctx context.Context, target utils.ValidatedTarget) domain.ModuleResult {
	result := scanner.NewResult(moduleName)
	client := utils.NewSafeHTTPClient(12 * time.Second)

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, target.CanonicalURL, nil)
	if err != nil {
		return scanner.Fail(moduleName, result.StartedAt, err)
	}
	req.Header.Set("User-Agent", "SecScan/1.0")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode == http.StatusMethodNotAllowed {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, target.CanonicalURL, nil)
		if err != nil {
			return scanner.Fail(moduleName, result.StartedAt, err)
		}
		req.Header.Set("User-Agent", "SecScan/1.0")
		resp, err = client.Do(req)
	}
	if err != nil {
		return scanner.Fail(moduleName, result.StartedAt, err)
	}
	defer resp.Body.Close()

	headerResults := make(map[string]map[string]string, len(rules))
	present := 0
	for _, rule := range rules {
		value := resp.Header.Get(rule.Name)
		if value == "" {
			result.Findings = append(result.Findings, domain.Finding{
				ID:             "headers-missing-" + strings.ToLower(rule.Name),
				Title:          fmt.Sprintf("Missing %s header", rule.Name),
				Description:    "A recommended browser security header was not present in the response.",
				Severity:       severityForHeader(rule.Name),
				Recommendation: rule.Recommendation,
			})
			headerResults[rule.Name] = map[string]string{"status": "missing"}
			continue
		}
		present++
		headerResults[rule.Name] = map[string]string{
			"status": "present",
			"value":  value,
		}
	}

	score := int(float64(present) / float64(len(rules)) * 100)
	if target.Scheme == "http" {
		score -= 10
		result.Findings = append(result.Findings, domain.Finding{
			ID:             "headers-http-no-hsts",
			Title:          "Target was scanned over HTTP",
			Description:    "HSTS only protects browsers after HTTPS has been used successfully.",
			Severity:       "medium",
			Recommendation: "Redirect HTTP traffic to HTTPS and enable Strict-Transport-Security on HTTPS responses.",
		})
	}

	result.Metadata["status_code"] = resp.StatusCode
	result.Metadata["header_results"] = headerResults
	result.Metadata["header_grade"] = utils.ScoreToGrade(score)

	scanner.Complete(&result, score)
	return result
}

func severityForHeader(name string) string {
	switch name {
	case "Content-Security-Policy", "Strict-Transport-Security":
		return "medium"
	default:
		return "low"
	}
}
