"use client";

import { FormEvent, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { MODULES, startScan } from "@/lib/api";

export default function HomePage() {
  const router = useRouter();
  const allModuleIDs = useMemo(() => MODULES.map((module) => module.id), []);
  const [url, setURL] = useState("");
  const [selectedModules, setSelectedModules] = useState<string[]>(allModuleIDs);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  function toggleModule(moduleID: string) {
    setSelectedModules((current) =>
      current.includes(moduleID) ? current.filter((item) => item !== moduleID) : [...current, moduleID]
    );
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");

    if (selectedModules.length === 0) {
      setError("Select at least one scanner module.");
      return;
    }

    try {
      const parsed = new URL(url);
      if (parsed.protocol !== "http:" && parsed.protocol !== "https:") {
        throw new Error("Only http and https URLs are supported.");
      }
    } catch {
      setError("Enter a valid http or https URL.");
      return;
    }

    setLoading(true);
    try {
      const response = await startScan(url, selectedModules);
      router.push(`/scan/${response.scan_id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not start the scan.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="min-h-screen px-4 py-8 sm:px-6 lg:px-10">
      <div className="mx-auto flex max-w-6xl flex-col gap-8">
        <header className="flex flex-wrap items-end justify-between gap-4 border-b border-[var(--line)] pb-6">
          <div>
            <p className="text-sm font-semibold uppercase tracking-normal text-emerald-700">Web Security Scanner Dashboard</p>
            <h1 className="mt-2 text-4xl font-bold text-[var(--ink)] sm:text-5xl">SecScan</h1>
          </div>
          <div className="rounded-full border border-[var(--line)] bg-white px-4 py-2 text-sm font-semibold text-zinc-700 shadow-soft">
            Academic scanner
          </div>
        </header>

        <section className="grid gap-6 lg:grid-cols-[1.05fr_0.95fr]">
          <form onSubmit={handleSubmit} className="rounded-lg border border-[var(--line)] bg-white p-6 shadow-soft">
            <label htmlFor="target-url" className="text-sm font-semibold text-zinc-800">
              Target URL
            </label>
            <div className="mt-2 flex flex-col gap-3 sm:flex-row">
              <input
                id="target-url"
                type="url"
                value={url}
                onChange={(event) => setURL(event.target.value)}
                placeholder="https://example.com"
                className="min-h-12 flex-1 rounded-md border border-[var(--line)] bg-white px-4 text-base outline-none transition focus:border-emerald-500 focus:ring-4 focus:ring-emerald-100"
                required
              />
              <button
                type="submit"
                disabled={loading}
                className="min-h-12 rounded-md bg-emerald-700 px-5 font-semibold text-white transition hover:bg-emerald-800 disabled:cursor-not-allowed disabled:bg-zinc-400"
              >
                {loading ? "Starting..." : "Start scan"}
              </button>
            </div>

            {error ? <p className="mt-3 rounded-md bg-rose-50 p-3 text-sm text-rose-800">{error}</p> : null}

            <div className="mt-6 flex flex-wrap items-center justify-between gap-3">
              <h2 className="text-lg font-semibold">Scanner Modules</h2>
              <div className="flex gap-2">
                <button
                  type="button"
                  onClick={() => setSelectedModules(allModuleIDs)}
                  className="rounded-md border border-[var(--line)] px-3 py-2 text-sm font-semibold text-zinc-700 hover:bg-zinc-50"
                >
                  Select all
                </button>
                <button
                  type="button"
                  onClick={() => setSelectedModules([])}
                  className="rounded-md border border-[var(--line)] px-3 py-2 text-sm font-semibold text-zinc-700 hover:bg-zinc-50"
                >
                  Clear
                </button>
              </div>
            </div>

            <div className="mt-4 grid gap-3 sm:grid-cols-2">
              {MODULES.map((module) => {
                const checked = selectedModules.includes(module.id);
                return (
                  <label
                    key={module.id}
                    className={`flex cursor-pointer gap-3 rounded-md border p-4 transition ${
                      checked ? "border-emerald-500 bg-emerald-50" : "border-[var(--line)] bg-white hover:bg-zinc-50"
                    }`}
                  >
                    <input
                      type="checkbox"
                      checked={checked}
                      onChange={() => toggleModule(module.id)}
                      className="mt-1 h-4 w-4 rounded border-zinc-300 text-emerald-700 focus:ring-emerald-600"
                    />
                    <span>
                      <span className="block font-semibold">{module.label}</span>
                      <span className="block text-sm text-[var(--muted)]">{module.detail}</span>
                    </span>
                  </label>
                );
              })}
            </div>
          </form>

          <aside className="rounded-lg border border-[var(--line)] bg-[#17201c] p-6 text-white shadow-soft">
            <h2 className="text-2xl font-bold">Security-first defaults</h2>
            <div className="mt-5 grid gap-4">
              {[
                ["SSRF guard", "Blocks localhost, private IP ranges, link-local addresses, and unsafe redirects."],
                ["Parallel modules", "Runs selected scanners concurrently and streams progress with Server-Sent Events."],
                ["Portable reports", "Generates a PDF report after completion with findings and recommendations."]
              ].map(([title, text]) => (
                <div key={title} className="border-t border-white/15 pt-4">
                  <h3 className="font-semibold text-emerald-200">{title}</h3>
                  <p className="mt-1 text-sm leading-6 text-zinc-200">{text}</p>
                </div>
              ))}
            </div>
          </aside>
        </section>
      </div>
    </main>
  );
}
