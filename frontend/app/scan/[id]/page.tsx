"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import { ModuleResultCard } from "@/components/ModuleResultCard";
import { RiskRadar } from "@/components/RiskRadar";
import { SeverityBadge } from "@/components/SeverityBadge";
import { StatusBadge } from "@/components/StatusBadge";
import { fetchScan, reportURL, streamURL } from "@/lib/api";
import type { ProgressEventPayload, Scan } from "@/lib/types";

const eventNames = ["scan_started", "module_started", "module_progress", "module_completed", "scan_completed", "scan_failed"];

export default function ScanResultPage() {
  const params = useParams<{ id: string }>();
  const scanID = params.id;
  const [scan, setScan] = useState<Scan | null>(null);
  const [progress, setProgress] = useState(0);
  const [events, setEvents] = useState<ProgressEventPayload[]>([]);
  const [error, setError] = useState("");
  const [streamConnected, setStreamConnected] = useState(false);

  useEffect(() => {
    let closed = false;

    async function load() {
      try {
        const current = await fetchScan(scanID);
        if (!closed) {
          setScan(current);
          setProgress(computeProgress(current));
        }
      } catch (err) {
        if (!closed) {
          setError(err instanceof Error ? err.message : "Could not load scan.");
        }
      }
    }

    load();

    const source = new EventSource(streamURL(scanID));
    setStreamConnected(true);

    const handleSnapshot = (event: Event) => {
      const message = event as MessageEvent<string>;
      const snapshot = JSON.parse(message.data) as Scan;
      setScan(snapshot);
      setProgress(computeProgress(snapshot));
    };

    const handleProgress = async (event: Event) => {
      const message = event as MessageEvent<string>;
      const payload = JSON.parse(message.data) as ProgressEventPayload;
      setEvents((current) => [payload, ...current].slice(0, 8));
      if (typeof payload.progress === "number") {
        setProgress(payload.progress);
      }
      if (payload.type === "module_completed" || payload.type === "scan_completed" || payload.type === "scan_failed") {
        try {
          const current = await fetchScan(scanID);
          setScan(current);
          setProgress(computeProgress(current));
        } catch {
          setError("Live update arrived, but the latest scan result could not be loaded.");
        }
      }
      if (payload.type === "scan_completed" || payload.type === "scan_failed") {
        source.close();
        setStreamConnected(false);
      }
    };

    source.addEventListener("snapshot", handleSnapshot);
    for (const eventName of eventNames) {
      source.addEventListener(eventName, handleProgress);
    }
    source.onerror = () => {
      setStreamConnected(false);
    };

    return () => {
      closed = true;
      source.close();
      setStreamConnected(false);
    };
  }, [scanID]);

  const orderedResults = useMemo(() => {
    if (!scan) {
      return [];
    }
    return scan.modules.map((module) => scan.results[module]).filter(Boolean);
  }, [scan]);

  if (error && !scan) {
    return (
      <main className="min-h-screen px-4 py-8">
        <div className="mx-auto max-w-3xl rounded-lg border border-rose-200 bg-white p-6 text-rose-800 shadow-soft">
          {error}
          <div className="mt-4">
            <Link href="/" className="font-semibold text-emerald-700">
              Back to scanner
            </Link>
          </div>
        </div>
      </main>
    );
  }

  if (!scan) {
    return (
      <main className="min-h-screen px-4 py-8">
        <div className="mx-auto max-w-3xl rounded-lg border border-[var(--line)] bg-white p-6 shadow-soft">Loading scan...</div>
      </main>
    );
  }

  return (
    <main className="min-h-screen px-4 py-8 sm:px-6 lg:px-10">
      <div className="mx-auto flex max-w-7xl flex-col gap-6">
        <header className="flex flex-wrap items-start justify-between gap-4 border-b border-[var(--line)] pb-6">
          <div>
            <Link href="/" className="text-sm font-semibold text-emerald-700 hover:text-emerald-900">
              Back to scanner
            </Link>
            <div className="mt-3 flex flex-wrap items-center gap-3">
              <h1 className="text-3xl font-bold sm:text-4xl">Scan Result</h1>
              <StatusBadge status={scan.status} />
            </div>
            <p className="mt-2 max-w-3xl break-all text-sm text-[var(--muted)]">{scan.canonical_url}</p>
          </div>
          <a
            href={reportURL(scan.id)}
            className={`rounded-md px-4 py-3 text-sm font-semibold text-white ${
              scan.status === "completed" ? "bg-emerald-700 hover:bg-emerald-800" : "pointer-events-none bg-zinc-400"
            }`}
          >
            Export PDF
          </a>
        </header>

        <section className="grid gap-4 md:grid-cols-4">
          <SummaryTile label="Grade" value={scan.overall_grade || "Pending"} detail={`${scan.overall_score || 0}/100`} strong />
          <SummaryTile label="Progress" value={`${progress}%`} detail={streamConnected ? "Live stream connected" : "Stream idle"} />
          <SummaryTile label="Modules" value={`${Object.keys(scan.results).length}/${scan.modules.length}`} detail="Completed or failed" />
          <SummaryTile label="Findings" value={String(scan.findings.length)} detail="Across selected modules" />
        </section>

        <div className="h-3 overflow-hidden rounded-full bg-zinc-200">
          <div className="h-full bg-emerald-600 transition-all duration-500" style={{ width: `${Math.min(progress, 100)}%` }} />
        </div>

        <section className="grid gap-6 lg:grid-cols-[0.95fr_1.05fr]">
          <RiskRadar modules={scan.modules} results={scan.results} />

          <div className="rounded-lg border border-[var(--line)] bg-white p-5 shadow-soft">
            <h2 className="text-lg font-semibold">Risk Highlights</h2>
            <div className="mt-4 flex flex-wrap gap-2">
              {Object.entries(scan.summary.findings_by_severity || {}).map(([severity, count]) => (
                <div key={severity} className="flex items-center gap-2 rounded-md border border-[var(--line)] px-3 py-2">
                  <SeverityBadge severity={severity} />
                  <span className="text-sm font-semibold">{count}</span>
                </div>
              ))}
            </div>
            <dl className="mt-5 grid gap-3 text-sm sm:grid-cols-2">
              <InfoItem label="Scan ID" value={scan.id} />
              <InfoItem label="Created" value={formatDate(scan.created_at)} />
              <InfoItem label="Started" value={scan.started_at ? formatDate(scan.started_at) : "Pending"} />
              <InfoItem label="Completed" value={scan.completed_at ? formatDate(scan.completed_at) : "Pending"} />
            </dl>
            {scan.error ? <p className="mt-4 rounded-md bg-rose-50 p-3 text-sm text-rose-800">{scan.error}</p> : null}
          </div>
        </section>

        <section className="rounded-lg border border-[var(--line)] bg-white p-5 shadow-soft">
          <h2 className="text-lg font-semibold">Live Events</h2>
          <div className="mt-4 space-y-2">
            {events.length === 0 ? (
              <p className="text-sm text-[var(--muted)]">Waiting for scanner events...</p>
            ) : (
              events.map((event) => (
                <div key={`${event.type}-${event.timestamp}-${event.module || "scan"}`} className="flex flex-wrap items-center gap-2 text-sm">
                  <span className="rounded-md bg-zinc-100 px-2 py-1 font-semibold text-zinc-700">{event.type}</span>
                  {event.module ? <span className="font-semibold">{event.module}</span> : null}
                  <span className="text-[var(--muted)]">{event.message}</span>
                </div>
              ))
            )}
          </div>
        </section>

        <section className="grid gap-5">
          {orderedResults.length === 0 ? (
            <div className="rounded-lg border border-[var(--line)] bg-white p-5 shadow-soft">Modules are still running.</div>
          ) : (
            orderedResults.map((result) => <ModuleResultCard key={result.name} result={result} />)
          )}
        </section>
      </div>
    </main>
  );
}

function SummaryTile({ label, value, detail, strong = false }: { label: string; value: string; detail: string; strong?: boolean }) {
  return (
    <div className="rounded-lg border border-[var(--line)] bg-white p-5 shadow-soft">
      <p className="text-sm font-semibold text-[var(--muted)]">{label}</p>
      <p className={`mt-2 font-bold ${strong ? "text-4xl text-emerald-700" : "text-2xl text-[var(--ink)]"}`}>{value}</p>
      <p className="mt-1 text-sm text-[var(--muted)]">{detail}</p>
    </div>
  );
}

function InfoItem({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-md bg-zinc-50 p-3">
      <dt className="font-semibold text-zinc-700">{label}</dt>
      <dd className="mt-1 break-words text-zinc-600">{value}</dd>
    </div>
  );
}

function computeProgress(scan: Scan): number {
  if (scan.status === "completed" || scan.status === "failed") {
    return 100;
  }
  if (scan.modules.length === 0) {
    return 0;
  }
  return Math.round((Object.keys(scan.results).length / scan.modules.length) * 100);
}

function formatDate(value: string): string {
  return new Intl.DateTimeFormat(undefined, {
    dateStyle: "medium",
    timeStyle: "short"
  }).format(new Date(value));
}
