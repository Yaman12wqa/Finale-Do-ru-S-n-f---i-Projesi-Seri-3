import type { Scan, StartScanResponse } from "./types";

export const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

export const MODULES = [
  { id: "ports", label: "Ports", detail: "Common TCP exposure" },
  { id: "headers", label: "Headers", detail: "Browser security headers" },
  { id: "tls", label: "TLS", detail: "Certificate and protocol checks" },
  { id: "fuzz", label: "Fuzz", detail: "Small safe path wordlist" },
  { id: "xss", label: "XSS", detail: "Reflected input heuristics" },
  { id: "sqli", label: "SQLi", detail: "Database error heuristics" },
  { id: "cve", label: "CVE", detail: "Technology advisory mapping" }
];

export async function startScan(url: string, modules: string[]): Promise<StartScanResponse> {
  const response = await fetch(`${API_BASE_URL}/api/scan`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ url, modules })
  });

  if (!response.ok) {
    throw new Error(await readError(response));
  }

  return response.json();
}

export async function fetchScan(scanID: string): Promise<Scan> {
  const response = await fetch(`${API_BASE_URL}/api/scan/${scanID}`, { cache: "no-store" });
  if (!response.ok) {
    throw new Error(await readError(response));
  }
  return response.json();
}

export function streamURL(scanID: string): string {
  return `${API_BASE_URL}/api/scan/${scanID}/stream`;
}

export function reportURL(scanID: string): string {
  return `${API_BASE_URL}/api/scan/${scanID}/report.pdf`;
}

async function readError(response: Response): Promise<string> {
  try {
    const data = await response.json();
    return data?.error?.message || `Request failed with status ${response.status}`;
  } catch {
    return `Request failed with status ${response.status}`;
  }
}
