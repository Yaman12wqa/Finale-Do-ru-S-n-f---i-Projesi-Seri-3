package fuzz

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/secscan/secscan/backend/internal/domain"
	"github.com/secscan/secscan/backend/internal/scanner"
	"github.com/secscan/secscan/backend/internal/utils"
)

const moduleName = "fuzz"

type Scanner struct{}

type discoveredPath struct {
	Path       string `json:"path"`
	StatusCode int    `json:"status_code"`
	Location   string `json:"location,omitempty"`
}

var wordlist = []string{
	"robots.txt",
	"sitemap.xml",
	".well-known/security.txt",
	"admin",
	"login",
	"dashboard",
	"api",
	"backup",
	"backup.zip",
	".env",
}

func New() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Name() string {
	return moduleName
}

func (s *Scanner) Scan(ctx context.Context, target utils.ValidatedTarget) domain.ModuleResult {
	result := scanner.NewResult(moduleName)
	client := utils.NewSafeHTTPClient(15 * time.Second)

	discovered := make([]discoveredPath, 0)
	score := 100
	for _, word := range wordlist {
		select {
		case <-ctx.Done():
			return scanner.Fail(moduleName, result.StartedAt, ctx.Err())
		default:
		}

		targetURL := pathURL(target.URL, word)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", "SecScan/1.0")

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()

		if isInterestingStatus(resp.StatusCode) {
			item := discoveredPath{
				Path:       "/" + strings.TrimLeft(word, "/"),
				StatusCode: resp.StatusCode,
				Location:   resp.Header.Get("Location"),
			}
			discovered = append(discovered, item)
			severity := "info"
			if word == ".env" || strings.Contains(word, "backup") {
				severity = "high"
				score -= 25
			}
			result.Findings = append(result.Findings, domain.Finding{
				ID:             "fuzz-" + strings.ReplaceAll(word, "/", "-"),
				Title:          "Interesting path discovered",
				Description:    "A lightweight path check returned a status code worth manual review.",
				Severity:       severity,
				Evidence:       fmt.Sprintf("%s returned HTTP %d", item.Path, item.StatusCode),
				Recommendation: "Confirm whether this path should be publicly reachable and remove sensitive files from web roots.",
			})
		}
	}

	result.Metadata["wordlist_size"] = len(wordlist)
	result.Metadata["discovered_paths"] = discovered
	if len(discovered) == 0 {
		result.Findings = append(result.Findings, domain.Finding{
			ID:             "fuzz-no-interesting-paths",
			Title:          "No interesting paths found",
			Description:    "The small built-in wordlist did not identify exposed administrative or metadata paths.",
			Severity:       "info",
			Recommendation: "Use authenticated and authorized testing for deeper content discovery when appropriate.",
		})
	}

	scanner.Complete(&result, score)
	return result
}

func pathURL(base *url.URL, word string) string {
	copied := *base
	copied.Path = "/" + strings.TrimLeft(word, "/")
	copied.RawQuery = ""
	copied.Fragment = ""
	return copied.String()
}

func isInterestingStatus(status int) bool {
	switch status {
	case http.StatusOK, http.StatusNoContent, http.StatusMovedPermanently, http.StatusFound,
		http.StatusTemporaryRedirect, http.StatusPermanentRedirect, http.StatusUnauthorized,
		http.StatusForbidden:
		return true
	default:
		return false
	}
}
