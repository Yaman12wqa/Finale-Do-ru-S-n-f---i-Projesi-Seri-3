package domain

import "time"

type ScanStatus string

const (
	ScanStatusPending   ScanStatus = "pending"
	ScanStatusRunning   ScanStatus = "running"
	ScanStatusCompleted ScanStatus = "completed"
	ScanStatusFailed    ScanStatus = "failed"
)

type ModuleStatus string

const (
	ModuleStatusRunning   ModuleStatus = "running"
	ModuleStatusCompleted ModuleStatus = "completed"
	ModuleStatusFailed    ModuleStatus = "failed"
)

const (
	EventScanStarted     = "scan_started"
	EventModuleStarted   = "module_started"
	EventModuleProgress  = "module_progress"
	EventModuleCompleted = "module_completed"
	EventScanCompleted   = "scan_completed"
	EventScanFailed      = "scan_failed"
)

type ScanRequest struct {
	URL     string   `json:"url" binding:"required"`
	Modules []string `json:"modules,omitempty"`
}

type StartScanResponse struct {
	ScanID          string     `json:"scan_id"`
	Status          ScanStatus `json:"status"`
	AcceptedURL     string     `json:"accepted_url"`
	SelectedModules []string   `json:"selected_modules"`
}

type Scan struct {
	ID           string                  `json:"id"`
	TargetURL    string                  `json:"target_url"`
	CanonicalURL string                  `json:"canonical_url"`
	Modules      []string                `json:"modules"`
	Status       ScanStatus              `json:"status"`
	OverallScore int                     `json:"overall_score"`
	OverallGrade string                  `json:"overall_grade"`
	Results      map[string]ModuleResult `json:"results"`
	Findings     []Finding               `json:"findings"`
	Summary      ReportSummary           `json:"summary"`
	CreatedAt    time.Time               `json:"created_at"`
	StartedAt    *time.Time              `json:"started_at,omitempty"`
	CompletedAt  *time.Time              `json:"completed_at,omitempty"`
	Error        string                  `json:"error,omitempty"`
}

type ModuleResult struct {
	Name        string         `json:"name"`
	Status      ModuleStatus   `json:"status"`
	Score       int            `json:"score"`
	Grade       string         `json:"grade"`
	Findings    []Finding      `json:"findings"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	StartedAt   time.Time      `json:"started_at"`
	CompletedAt time.Time      `json:"completed_at"`
	DurationMS  int64          `json:"duration_ms"`
	Error       string         `json:"error,omitempty"`
}

type Finding struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Severity       string `json:"severity"`
	Evidence       string `json:"evidence,omitempty"`
	Recommendation string `json:"recommendation"`
}

type ProgressEvent struct {
	Type      string        `json:"type"`
	ScanID    string        `json:"scan_id"`
	Module    string        `json:"module,omitempty"`
	Status    ScanStatus    `json:"status,omitempty"`
	Progress  int           `json:"progress"`
	Message   string        `json:"message"`
	Result    *ModuleResult `json:"result,omitempty"`
	Error     string        `json:"error,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
}

type ReportSummary struct {
	ScanID             string         `json:"scan_id,omitempty"`
	TargetURL          string         `json:"target_url,omitempty"`
	OverallGrade       string         `json:"overall_grade,omitempty"`
	OverallScore       int            `json:"overall_score,omitempty"`
	CompletedModules   int            `json:"completed_modules"`
	FailedModules      int            `json:"failed_modules"`
	FindingsBySeverity map[string]int `json:"findings_by_severity"`
}
