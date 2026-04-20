# SecScan - Web Security Scanner Dashboard

Student: YMAN ALSHEABE
School No: 23080410056

Academic use only. Run SecScan only against systems you own or have explicit written permission to test.

## Overview

SecScan is a university Web Programming project that provides a web dashboard for lightweight security scanning. A user submits a target URL, the backend validates it with SSRF protections, runs selected scanner modules in parallel, streams live progress through Server-Sent Events, and produces a scored JSON result plus a downloadable PDF report.

The project is intentionally realistic but explainable: the backend uses Go, Gin, clean service boundaries, a scanner registry, in-memory storage, safe timeouts, and a pluggable CVE adapter. The frontend uses Next.js App Router, TypeScript, and Tailwind CSS.

## Tech Stack

- Backend: Go 1.22+, Gin
- Frontend: Next.js 14, TypeScript, Tailwind CSS
- Reports: gofpdf
- Infra: Docker, Docker Compose
- CI/Security: GitHub Actions, Semgrep, Trivy

## Folder Structure

```text
.
|-- backend/
|   |-- Dockerfile
|   |-- go.mod
|   |-- main.go
|   `-- internal/
|       |-- api/
|       |-- config/
|       |-- domain/
|       |-- report/
|       |-- scanner/
|       |-- service/
|       |-- sse/
|       |-- storage/
|       `-- utils/
|-- frontend/
|   |-- Dockerfile
|   |-- app/
|   |-- components/
|   |-- lib/
|   `-- package.json
|-- gorev-01/ ... gorev-10/
|-- docs/
|-- .github/workflows/
|-- docker-compose.yml
`-- README.md
```

The `gorev-01` through `gorev-10` folders contain Kisim 1 mini-task reports, command notes, and screenshot placeholders. They are separated from the SecScan full-stack application.

## Setup

### Requirements

- Docker and Docker Compose
- Optional for local development without Docker:
  - Go 1.22+
  - Node.js 20+
  - npm

### Run With Docker Compose

```bash
docker compose up --build
```

Open:

- Frontend: http://localhost:3000
- Backend health check: http://localhost:8080/health

Stop:

```bash
docker compose down
```

## Environment Variables

Root `.env.example`:

```env
PORT=8080
GIN_MODE=release
FRONTEND_ORIGIN=http://localhost:3000
SCAN_TIMEOUT_SECONDS=60
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

Backend variables:

- `PORT`: backend HTTP port
- `GIN_MODE`: Gin runtime mode, usually `release`
- `FRONTEND_ORIGIN`: allowed CORS origin
- `SCAN_TIMEOUT_SECONDS`: maximum scan runtime

Frontend variables:

- `NEXT_PUBLIC_API_BASE_URL`: browser-visible backend URL

## Local Development

Backend:

```bash
cd backend
go mod download
go test ./...
go run .
```

Frontend:

```bash
cd frontend
npm ci
npm run dev
```

## API Endpoints

### `GET /health`

Returns:

```json
{ "status": "ok" }
```

### `POST /api/scan`

Starts a scan asynchronously.

Request:

```json
{
  "url": "https://example.com",
  "modules": ["ports", "headers", "tls", "fuzz", "xss", "sqli", "cve"]
}
```

`modules` is optional. If omitted, all modules run.

Response:

```json
{
  "scan_id": "generated-id",
  "status": "pending",
  "accepted_url": "https://example.com",
  "selected_modules": ["ports", "headers"]
}
```

### `GET /api/scan/:id`

Returns the full scan result as JSON.

### `GET /api/scan/:id/stream`

Streams live progress with Server-Sent Events:

- `scan_started`
- `module_started`
- `module_progress`
- `module_completed`
- `scan_completed`
- `scan_failed`

### `GET /api/scan/:id/report.pdf`

Returns a generated PDF report after the scan completes.

## Scanner Modules

- `ports`: safe TCP connect scan using a goroutine worker pool and common service labels.
- `headers`: checks CSP, HSTS, X-Frame-Options, X-Content-Type-Options, Referrer-Policy, Permissions-Policy, COOP, and CORP.
- `tls`: checks TLS support and certificate subject, issuer, validity, and verification warnings.
- `fuzz`: uses a small built-in wordlist for low-volume path discovery.
- `xss`: checks safe reflected input probes in query parameters and reports heuristic reflection only.
- `sqli`: checks safe SQL syntax probes for database error signatures and reports heuristic findings.
- `cve`: detects likely technologies and maps them through a local pluggable advisory adapter.

## SSRF Protection

The backend validates the target before scanning and also revalidates outbound HTTP redirects and dial targets. It only allows `http` and `https`, resolves DNS, and blocks localhost, loopback, private, link-local, multicast, unspecified, and local IPv6 ranges.

Blocked examples include:

- `127.0.0.0/8`
- `10.0.0.0/8`
- `172.16.0.0/12`
- `192.168.0.0/16`
- `169.254.0.0/16`
- `localhost`
- IPv6 loopback and local ranges such as `::1`, `fc00::/7`, and `fe80::/10`

The scanner is designed for authorized academic testing only. It is not a replacement for professional penetration testing.

## Screenshots

Add screenshots before submission:

- Home page: `docs/screenshots/home.png`
- Backend health check: `docs/screenshots/health.png`
- Live scan page: `docs/screenshots/scan-running.png`
- Completed report page: `docs/screenshots/scan-completed.png`
- PDF example: `docs/screenshots/report.png`

## CI and Security

GitHub workflows are included:

- `.github/workflows/ci.yml`: backend test/build and frontend lint/typecheck/build
- `.github/workflows/security.yml`: Semgrep rules and Trivy filesystem scan

## AI Usage Transparency

This project was generated with AI assistance and should be reviewed, tested, and understood by the student before submission. The student is responsible for verifying correctness, explaining the architecture, and ensuring the project follows course rules.

## Limitations

- Results are heuristic and must be manually verified.
- Storage is in-memory, so scan history is lost when the backend restarts.
- The CVE module uses a local mock adapter instead of a live NVD/OSV integration.
- The fuzzing module intentionally uses a small wordlist to avoid aggressive scanning.
- XSS and SQLi modules do not claim confirmed exploitation.

## Future Improvements

- Add persistent storage for scan history.
- Add authenticated scan ownership and user accounts.
- Replace the local CVE adapter with an NVD or OSV-backed adapter.
- Add configurable port and path wordlists.
- Add rate limiting and scan quotas for multi-user deployments.
