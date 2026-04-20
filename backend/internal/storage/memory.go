package storage

import (
	"sync"

	"github.com/secscan/secscan/backend/internal/domain"
)

type MemoryStore struct {
	mu    sync.RWMutex
	scans map[string]*domain.Scan
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{scans: make(map[string]*domain.Scan)}
}

func (s *MemoryStore) Create(scan domain.Scan) {
	s.mu.Lock()
	defer s.mu.Unlock()
	copied := cloneScan(scan)
	s.scans[scan.ID] = &copied
}

func (s *MemoryStore) Get(id string) (domain.Scan, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	scan, ok := s.scans[id]
	if !ok {
		return domain.Scan{}, false
	}
	return cloneScan(*scan), true
}

func (s *MemoryStore) Update(id string, update func(*domain.Scan)) (domain.Scan, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	scan, ok := s.scans[id]
	if !ok {
		return domain.Scan{}, false
	}
	update(scan)
	return cloneScan(*scan), true
}

func cloneScan(scan domain.Scan) domain.Scan {
	copied := scan
	copied.Modules = append([]string(nil), scan.Modules...)
	copied.Findings = append([]domain.Finding(nil), scan.Findings...)
	copied.Results = make(map[string]domain.ModuleResult, len(scan.Results))
	for key, value := range scan.Results {
		copied.Results[key] = cloneModuleResult(value)
	}
	if scan.Summary.FindingsBySeverity != nil {
		copied.Summary.FindingsBySeverity = make(map[string]int, len(scan.Summary.FindingsBySeverity))
		for key, value := range scan.Summary.FindingsBySeverity {
			copied.Summary.FindingsBySeverity[key] = value
		}
	}
	return copied
}

func cloneModuleResult(result domain.ModuleResult) domain.ModuleResult {
	copied := result
	copied.Findings = append([]domain.Finding(nil), result.Findings...)
	if result.Metadata != nil {
		copied.Metadata = make(map[string]any, len(result.Metadata))
		for key, value := range result.Metadata {
			copied.Metadata[key] = value
		}
	}
	return copied
}
