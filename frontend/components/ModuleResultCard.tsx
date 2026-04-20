import type { ModuleResult } from "@/lib/types";
import { SeverityBadge } from "./SeverityBadge";
import { StatusBadge } from "./StatusBadge";

export function ModuleResultCard({ result }: { result: ModuleResult }) {
  return (
    <article className="rounded-lg border border-[var(--line)] bg-white p-5 shadow-soft">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h3 className="text-lg font-semibold uppercase tracking-normal">{result.name}</h3>
          <p className="text-sm text-[var(--muted)]">{result.duration_ms} ms</p>
        </div>
        <div className="flex items-center gap-2">
          <StatusBadge status={result.status} />
          <span className="rounded-full bg-zinc-100 px-3 py-1 text-sm font-bold text-zinc-800">
            {result.grade} / {result.score}
          </span>
        </div>
      </div>

      {result.error ? <p className="mt-4 rounded-md bg-rose-50 p-3 text-sm text-rose-800">{result.error}</p> : null}

      <div className="mt-4 space-y-3">
        {result.findings.length === 0 ? (
          <p className="text-sm text-[var(--muted)]">No findings reported by this module.</p>
        ) : (
          result.findings.map((finding) => (
            <div key={finding.id} className="border-t border-[var(--line)] pt-3">
              <div className="flex flex-wrap items-center gap-2">
                <SeverityBadge severity={finding.severity} />
                <h4 className="font-semibold">{finding.title}</h4>
              </div>
              <p className="mt-2 text-sm text-zinc-700">{finding.description}</p>
              {finding.evidence ? <p className="mt-2 break-words text-sm text-zinc-600">Evidence: {finding.evidence}</p> : null}
              <p className="mt-2 text-sm text-zinc-700">Recommendation: {finding.recommendation}</p>
            </div>
          ))
        )}
      </div>

      {result.metadata ? (
        <details className="mt-4 rounded-md border border-[var(--line)] bg-zinc-50 p-3">
          <summary className="cursor-pointer text-sm font-semibold text-zinc-700">Raw module metadata</summary>
          <pre className="mt-3 max-h-72 overflow-auto whitespace-pre-wrap break-words text-xs text-zinc-700">
            {JSON.stringify(result.metadata, null, 2)}
          </pre>
        </details>
      ) : null}
    </article>
  );
}
