package xss

import (
	"context"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/secscan/secscan/backend/internal/domain"
	"github.com/secscan/secscan/backend/internal/scanner"
	"github.com/secscan/secscan/backend/internal/utils"
)

const moduleName = "xss"

type Scanner struct{}

func New() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Name() string {
	return moduleName
}

func (s *Scanner) Scan(ctx context.Context, target utils.ValidatedTarget) domain.ModuleResult {
	result := scanner.NewResult(moduleName)
	params := target.URL.Query()
	if len(params) == 0 {
		result.Findings = append(result.Findings, domain.Finding{
			ID:             "xss-no-query-params",
			Title:          "No query parameters available for reflection testing",
			Description:    "The submitted URL did not include query parameters, so reflected input checks were skipped.",
			Severity:       "info",
			Recommendation: "Test parameterized application routes during authorized dynamic testing.",
		})
		result.Metadata["tested_parameters"] = []string{}
		scanner.Complete(&result, 100)
		return result
	}

	client := utils.NewSafeHTTPClient(12 * time.Second)
	payload := "secscan_xss_probe_7b3c<>\"'"
	escapedPayload := html.EscapeString(payload)
	tested := make([]string, 0)
	reflections := 0

	for param := range params {
		if len(tested) >= 3 {
			break
		}
		select {
		case <-ctx.Done():
			return scanner.Fail(moduleName, result.StartedAt, ctx.Err())
		default:
		}

		tested = append(tested, param)
		testURL := withQueryParam(target.URL, param, payload)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, testURL, nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", "SecScan/1.0")

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		body, err := utils.ReadLimited(resp, 512*1024)
		if err != nil {
			continue
		}

		if containsReflection(body, payload, escapedPayload) {
			reflections++
			result.Findings = append(result.Findings, domain.Finding{
				ID:             "xss-reflection-" + param,
				Title:          "Reflected input observed",
				Description:    "A safe probe value was reflected in the response body. This is a heuristic signal, not a confirmed XSS exploit.",
				Severity:       "medium",
				Evidence:       fmt.Sprintf("Parameter %q reflected the SecScan probe", param),
				Recommendation: "Contextually encode output and validate input before rendering user-controlled values.",
			})
		}
	}

	result.Metadata["tested_parameters"] = tested
	result.Metadata["reflections"] = reflections

	score := 95
	if reflections > 0 {
		score = 70
	}
	scanner.Complete(&result, score)
	return result
}

func withQueryParam(base *url.URL, key, value string) string {
	copied := *base
	query := copied.Query()
	query.Set(key, value)
	copied.RawQuery = query.Encode()
	return copied.String()
}

func containsReflection(body, payload, escapedPayload string) bool {
	return strings.Contains(body, payload) || strings.Contains(body, escapedPayload)
}
