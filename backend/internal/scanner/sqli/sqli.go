package sqli

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

const moduleName = "sqli"

type Scanner struct{}

var payloads = []string{"'", "\"", "1' OR '1'='1"}

var errorSignatures = []string{
	"you have an error in your sql syntax",
	"warning: mysql",
	"unclosed quotation mark",
	"quoted string not properly terminated",
	"postgresql",
	"sqlite error",
	"sqlstate",
	"ora-",
	"odbc",
	"syntax error at or near",
}

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
			ID:             "sqli-no-query-params",
			Title:          "No query parameters available for SQLi heuristics",
			Description:    "The submitted URL did not include query parameters, so lightweight SQL injection checks were skipped.",
			Severity:       "info",
			Recommendation: "Test parameterized routes with permission and confirm findings manually.",
		})
		result.Metadata["tested_parameters"] = []string{}
		scanner.Complete(&result, 100)
		return result
	}

	client := utils.NewSafeHTTPClient(12 * time.Second)
	tested := make([]string, 0)
	signals := 0

	for param := range params {
		if len(tested) >= 3 {
			break
		}
		tested = append(tested, param)
		paramHadSignal := false
		for _, payload := range payloads {
			select {
			case <-ctx.Done():
				return scanner.Fail(moduleName, result.StartedAt, ctx.Err())
			default:
			}

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

			signature := findSQLSignature(body)
			if signature != "" {
				signals++
				paramHadSignal = true
				result.Findings = append(result.Findings, domain.Finding{
					ID:             "sqli-error-" + param,
					Title:          "SQL error signature observed",
					Description:    "A safe SQL syntax probe produced a response containing a database error signature. This is a heuristic signal, not confirmed exploitation.",
					Severity:       "high",
					Evidence:       fmt.Sprintf("Parameter %q, signature %q", param, signature),
					Recommendation: "Use parameterized queries, avoid exposing database errors, and manually verify the affected route.",
				})
				break
			}
		}
		if paramHadSignal {
			continue
		}
	}

	result.Metadata["tested_parameters"] = tested
	result.Metadata["heuristic_signals"] = signals

	score := 95
	if signals > 0 {
		score = 60
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

func findSQLSignature(body string) string {
	lower := strings.ToLower(body)
	for _, signature := range errorSignatures {
		if strings.Contains(lower, signature) {
			return signature
		}
	}
	return ""
}
