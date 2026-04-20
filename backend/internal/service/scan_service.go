package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/secscan/secscan/backend/internal/domain"
	"github.com/secscan/secscan/backend/internal/scanner"
	"github.com/secscan/secscan/backend/internal/scanner/registry"
	"github.com/secscan/secscan/backend/internal/sse"
	"github.com/secscan/secscan/backend/internal/storage"
	"github.com/secscan/secscan/backend/internal/utils"
)

type ScanService struct {
	store       *storage.MemoryStore
	broker      *sse.Broker
	registry    *registry.Registry
	scanTimeout time.Duration
}

func NewScanService(store *storage.MemoryStore, broker *sse.Broker, registry *registry.Registry, scanTimeout time.Duration) *ScanService {
	return &ScanService{
		store:       store,
		broker:      broker,
		registry:    registry,
		scanTimeout: scanTimeout,
	}
}

func (s *ScanService) StartScan(ctx context.Context, request domain.ScanRequest) (domain.StartScanResponse, error) {
	target, err := utils.ValidateTargetURL(ctx, request.URL)
	if err != nil {
		return domain.StartScanResponse{}, fmt.Errorf("target rejected by SSRF guard: %w", err)
	}

	selected, selectedNames, err := s.registry.Select(request.Modules)
	if err != nil {
		return domain.StartScanResponse{}, err
	}

	id, err := newScanID()
	if err != nil {
		return domain.StartScanResponse{}, err
	}

	scan := domain.Scan{
		ID:           id,
		TargetURL:    request.URL,
		CanonicalURL: target.CanonicalURL,
		Modules:      selectedNames,
		Status:       domain.ScanStatusPending,
		OverallScore: 0,
		OverallGrade: "F",
		Results:      make(map[string]domain.ModuleResult),
		Findings:     []domain.Finding{},
		Summary: domain.ReportSummary{
			ScanID:             id,
			TargetURL:          target.CanonicalURL,
			FindingsBySeverity: map[string]int{},
		},
		CreatedAt: time.Now().UTC(),
	}
	s.store.Create(scan)

	go s.runScan(id, target, selected)

	return domain.StartScanResponse{
		ScanID:          id,
		Status:          domain.ScanStatusPending,
		AcceptedURL:     target.CanonicalURL,
		SelectedModules: selectedNames,
	}, nil
}

func (s *ScanService) GetScan(id string) (domain.Scan, bool) {
	return s.store.Get(id)
}

func (s *ScanService) Subscribe(id string) (<-chan domain.ProgressEvent, func()) {
	return s.broker.Subscribe(id)
}

func (s *ScanService) runScan(scanID string, target utils.ValidatedTarget, modules []scanner.Scanner) {
	startedAt := time.Now().UTC()
	s.store.Update(scanID, func(scan *domain.Scan) {
		scan.Status = domain.ScanStatusRunning
		scan.StartedAt = &startedAt
	})
	s.publish(domain.ProgressEvent{
		Type:      domain.EventScanStarted,
		ScanID:    scanID,
		Status:    domain.ScanStatusRunning,
		Progress:  0,
		Message:   "scan started",
		Timestamp: startedAt,
	})

	ctx, cancel := context.WithTimeout(context.Background(), s.scanTimeout)
	defer cancel()

	resultCh := make(chan domain.ModuleResult, len(modules))
	var wg sync.WaitGroup

	for _, module := range modules {
		wg.Add(1)
		go func(module scanner.Scanner) {
			defer wg.Done()
			s.publish(domain.ProgressEvent{
				Type:      domain.EventModuleStarted,
				ScanID:    scanID,
				Module:    module.Name(),
				Status:    domain.ScanStatusRunning,
				Message:   fmt.Sprintf("%s module started", module.Name()),
				Timestamp: time.Now().UTC(),
			})
			resultCh <- module.Scan(ctx, target)
		}(module)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	completed := 0
	total := len(modules)
	for result := range resultCh {
		completed++
		progress := int(float64(completed) / float64(total) * 100)
		stored, _ := s.store.Update(scanID, func(scan *domain.Scan) {
			scan.Results[result.Name] = result
			aggregateScan(scan)
		})

		resultCopy := result
		s.publish(domain.ProgressEvent{
			Type:      domain.EventModuleCompleted,
			ScanID:    scanID,
			Module:    result.Name,
			Status:    stored.Status,
			Progress:  progress,
			Message:   fmt.Sprintf("%s module completed", result.Name),
			Result:    &resultCopy,
			Timestamp: time.Now().UTC(),
		})
		s.publish(domain.ProgressEvent{
			Type:      domain.EventModuleProgress,
			ScanID:    scanID,
			Module:    result.Name,
			Status:    stored.Status,
			Progress:  progress,
			Message:   fmt.Sprintf("%d of %d modules completed", completed, total),
			Timestamp: time.Now().UTC(),
		})
	}

	completedAt := time.Now().UTC()
	finalScan, _ := s.store.Update(scanID, func(scan *domain.Scan) {
		aggregateScan(scan)
		scan.CompletedAt = &completedAt
		if scan.Summary.CompletedModules == 0 {
			scan.Status = domain.ScanStatusFailed
			scan.Error = "all scanner modules failed"
		} else {
			scan.Status = domain.ScanStatusCompleted
		}
	})

	eventType := domain.EventScanCompleted
	message := "scan completed"
	if finalScan.Status == domain.ScanStatusFailed {
		eventType = domain.EventScanFailed
		message = finalScan.Error
	}
	s.publish(domain.ProgressEvent{
		Type:      eventType,
		ScanID:    scanID,
		Status:    finalScan.Status,
		Progress:  100,
		Message:   message,
		Error:     finalScan.Error,
		Timestamp: completedAt,
	})
}

func (s *ScanService) publish(event domain.ProgressEvent) {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	s.broker.Publish(event)
}

func aggregateScan(scan *domain.Scan) {
	if len(scan.Modules) == 0 {
		return
	}

	totalScore := 0
	completedModules := 0
	failedModules := 0
	findings := make([]domain.Finding, 0)
	bySeverity := map[string]int{
		"critical": 0,
		"high":     0,
		"medium":   0,
		"low":      0,
		"info":     0,
	}

	for _, moduleName := range scan.Modules {
		result, exists := scan.Results[moduleName]
		if !exists {
			continue
		}
		totalScore += result.Score
		if result.Status == domain.ModuleStatusFailed {
			failedModules++
		} else {
			completedModules++
		}
		for _, finding := range result.Findings {
			findings = append(findings, finding)
			if _, exists := bySeverity[finding.Severity]; !exists {
				bySeverity[finding.Severity] = 0
			}
			bySeverity[finding.Severity]++
		}
	}

	scoredModules := completedModules + failedModules
	if scoredModules > 0 {
		scan.OverallScore = utils.ClampScore(totalScore / scoredModules)
		scan.OverallGrade = utils.ScoreToGrade(scan.OverallScore)
	}
	scan.Findings = findings
	scan.Summary = domain.ReportSummary{
		ScanID:             scan.ID,
		TargetURL:          scan.CanonicalURL,
		OverallGrade:       scan.OverallGrade,
		OverallScore:       scan.OverallScore,
		CompletedModules:   completedModules,
		FailedModules:      failedModules,
		FindingsBySeverity: bySeverity,
	}
}

func newScanID() (string, error) {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes[:]), nil
}
