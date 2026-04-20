# SecScan Requirement Check

Student: YMAN ALSHEABE
School No: 23080410056

This file compares the generated SecScan project with the assignment screenshots.

## Kisim 2 - Full Stack Security Scanner

| Requirement | Status | Evidence |
|---|---:|---|
| Backend in Go + Gin | Done | `backend/go.mod`, `backend/main.go` |
| Frontend in Next.js 14 + TypeScript + Tailwind | Done | `frontend/package.json`, `frontend/app`, `frontend/tailwind.config.ts` |
| Docker Compose local run | Done | `docker-compose.yml`, backend and frontend Dockerfiles |
| Separate backend/frontend folders | Done | `backend/`, `frontend/` |
| `GET /health` | Done | `backend/internal/api/routes/routes.go` |
| `POST /api/scan` | Done | asynchronous scan start |
| `GET /api/scan/:id` | Done | full JSON result |
| `GET /api/scan/:id/stream` | Done | SSE stream |
| `GET /api/scan/:id/report.pdf` | Done | PDF report generation |
| Scanner registry architecture | Done | `backend/internal/scanner/registry` |
| 7 scanner modules | Done | ports, headers, tls, fuzz, xss, sqli, cve |
| Parallel module execution | Done | goroutines in `service/scan_service.go` |
| Scan lifecycle | Done | pending, running, completed, failed |
| SSRF guard | Done | DNS resolution and blocked CIDR checks in `utils/ssrf.go` |
| Safe HTTP timeouts and redirects | Done | `utils/http.go` |
| CORS config | Done | `api/middleware/cors.go` |
| PDF report | Done | `report/pdf.go` |
| README placeholders in first 10 lines | Done | README lines 3-4 contain student data |
| CI workflow | Done | `.github/workflows/ci.yml` |
| Security workflow | Done | `.github/workflows/security.yml` |

## Kisim 1 - 10 Mini Tasks

The repository now has `gorev-01` through `gorev-10` folders. Each folder contains a `rapor.md` with:

1. task number and title,
2. numbered steps,
3. encountered errors / debugging notes,
4. result and screenshot placeholders,
5. three learning bullets.

Some tasks require manual evidence from external tools or websites. Those cannot be fabricated:

- Juice Shop / DVWA screenshots,
- Google OAuth client configuration,
- OWASP ZAP HTML report,
- GitHub Actions run screenshots,
- Syft / Trivy output,
- securityheaders.com permalink.

## Current Local Verification

The following checks were run successfully before this file was created:

```text
go test ./...
go build ./...
npm run lint
npm run typecheck
npm run build
docker compose config
docker compose build
docker compose up -d
GET http://localhost:8080/health -> status ok
GET http://localhost:3000 -> HTTP 200
headers-only scan against https://example.com -> completed, PDF HTTP 200
```

## Remaining Manual Submission Work

- Create a GitHub repository/branch as required by the teacher.
- Commit changes regularly instead of one final commit.
- Capture screenshots into each `gorev-XX/ekran-goruntuleri/` folder.
- Fill the "Kanitim" and "Karsilastigim Hatalar" sections with your own exact results.
- Add final demo screenshots for SecScan under the README screenshots section.
