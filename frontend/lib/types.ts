export type ScanStatus = "pending" | "running" | "completed" | "failed";
export type ModuleStatus = "running" | "completed" | "failed";
export type Severity = "critical" | "high" | "medium" | "low" | "info" | string;

export interface Finding {
  id: string;
  title: string;
  description: string;
  severity: Severity;
  evidence?: string;
  recommendation: string;
}

export interface ModuleResult {
  name: string;
  status: ModuleStatus;
  score: number;
  grade: string;
  findings: Finding[];
  metadata?: Record<string, unknown>;
  started_at: string;
  completed_at: string;
  duration_ms: number;
  error?: string;
}

export interface ReportSummary {
  scan_id?: string;
  target_url?: string;
  overall_grade?: string;
  overall_score?: number;
  completed_modules: number;
  failed_modules: number;
  findings_by_severity: Record<string, number>;
}

export interface Scan {
  id: string;
  target_url: string;
  canonical_url: string;
  modules: string[];
  status: ScanStatus;
  overall_score: number;
  overall_grade: string;
  results: Record<string, ModuleResult>;
  findings: Finding[];
  summary: ReportSummary;
  created_at: string;
  started_at?: string;
  completed_at?: string;
  error?: string;
}

export interface StartScanResponse {
  scan_id: string;
  status: ScanStatus;
  accepted_url: string;
  selected_modules: string[];
}

export interface ProgressEventPayload {
  type: string;
  scan_id: string;
  module?: string;
  status?: ScanStatus;
  progress: number;
  message: string;
  result?: ModuleResult;
  error?: string;
  timestamp: string;
}
