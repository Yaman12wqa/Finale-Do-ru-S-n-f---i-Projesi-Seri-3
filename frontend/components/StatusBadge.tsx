import type { ScanStatus, ModuleStatus } from "@/lib/types";

const statusClasses: Record<string, string> = {
  pending: "bg-amber-100 text-amber-800 ring-amber-200",
  running: "bg-sky-100 text-sky-800 ring-sky-200",
  completed: "bg-emerald-100 text-emerald-800 ring-emerald-200",
  failed: "bg-rose-100 text-rose-800 ring-rose-200"
};

export function StatusBadge({ status }: { status: ScanStatus | ModuleStatus }) {
  return (
    <span className={`inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold ring-1 ${statusClasses[status] || statusClasses.pending}`}>
      {status}
    </span>
  );
}
