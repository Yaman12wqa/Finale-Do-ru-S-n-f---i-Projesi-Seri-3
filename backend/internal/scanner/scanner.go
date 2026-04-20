package scanner

import (
	"context"
	"time"

	"github.com/secscan/secscan/backend/internal/domain"
	"github.com/secscan/secscan/backend/internal/utils"
)

type Scanner interface {
	Name() string
	Scan(ctx context.Context, target utils.ValidatedTarget) domain.ModuleResult
}

func NewResult(name string) domain.ModuleResult {
	return domain.ModuleResult{
		Name:      name,
		Status:    domain.ModuleStatusRunning,
		Findings:  []domain.Finding{},
		Metadata:  map[string]any{},
		StartedAt: time.Now().UTC(),
	}
}

func Complete(result *domain.ModuleResult, score int) {
	now := time.Now().UTC()
	result.Status = domain.ModuleStatusCompleted
	result.Score = utils.ClampScore(score)
	result.Grade = utils.ScoreToGrade(result.Score)
	result.CompletedAt = now
	result.DurationMS = now.Sub(result.StartedAt).Milliseconds()
}

func Fail(name string, started time.Time, err error) domain.ModuleResult {
	now := time.Now().UTC()
	return domain.ModuleResult{
		Name:        name,
		Status:      domain.ModuleStatusFailed,
		Score:       0,
		Grade:       "F",
		Findings:    []domain.Finding{},
		Metadata:    map[string]any{},
		StartedAt:   started,
		CompletedAt: now,
		DurationMS:  now.Sub(started).Milliseconds(),
		Error:       err.Error(),
	}
}
