import type { Severity } from "@/lib/types";

const severityClasses: Record<string, string> = {
  critical: "bg-rose-950 text-white",
  high: "bg-rose-100 text-rose-800 ring-rose-200",
  medium: "bg-amber-100 text-amber-800 ring-amber-200",
  low: "bg-sky-100 text-sky-800 ring-sky-200",
  info: "bg-zinc-100 text-zinc-700 ring-zinc-200"
};

export function SeverityBadge({ severity }: { severity: Severity }) {
  const normalized = String(severity).toLowerCase();
  return (
    <span className={`inline-flex items-center rounded-full px-2.5 py-1 text-xs font-semibold ring-1 ${severityClasses[normalized] || severityClasses.info}`}>
      {normalized}
    </span>
  );
}
