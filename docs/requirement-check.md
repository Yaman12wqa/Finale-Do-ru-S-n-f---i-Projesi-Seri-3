# SecScan Requirement Check

This file compares the generated SecScan project with the assignment screenshots.

The status column only records whether an implementation is present in the repository. It does not prove that the project owner wrote, tested, or can explain that implementation.

## Kisim 2 - Full Stack Security Scanner

| Requirement | Status | Evidence |
|---|---:|---|
| Backend in Go + Gin | Present | `backend/go.mod`, `backend/main.go` |
| Frontend in Next.js 14 + TypeScript + Tailwind | Present | `frontend/package.json`, `frontend/app`, `frontend/tailwind.config.ts` |
| Docker Compose files | Present | `docker-compose.yml`, backend and frontend Dockerfiles |
| Separate backend/frontend folders | Present | `backend/`, `frontend/` |
| `GET /health` | Present | `backend/internal/api/routes/routes.go` |
| `POST /api/scan` | Present | asynchronous scan start |
| `GET /api/scan/:id` | Present | full JSON result |
| `GET /api/scan/:id/stream` | Present | SSE stream |
| `GET /api/scan/:id/report.pdf` | Present | PDF report generation |
| Scanner registry architecture | Present | `backend/internal/scanner/registry` |
| 7 scanner modules | Present | ports, headers, tls, fuzz, xss, sqli, cve |
| Parallel module execution | Present | goroutines in `service/scan_service.go` |
| Scan lifecycle | Present | pending, running, completed, failed |
| SSRF guard | Present | DNS resolution and blocked CIDR checks in `utils/ssrf.go` |
| HTTP timeout and redirect logic | Present | `utils/http.go` |
| CORS config | Present | `api/middleware/cors.go` |
| PDF report code | Present | `report/pdf.go` |
| CI workflow file | Present | `.github/workflows/ci.yml` |
| Security workflow file | Present | `.github/workflows/security.yml` |

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

## Historical Verification Note

An earlier version of this document claimed that the following checks passed. Treat that as unconfirmed until the project owner reruns the commands, records the current results, and can explain what each check covers:

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
