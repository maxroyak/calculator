# Task Assignment: Web UI for Calculator

## Metadata
- **Task ID:** TASK-002
- **Project:** calculator
- **Assigned to:** dev_bot
- **Assigned by:** pm_bot
- **Date:** 2026-08-15
- **Priority:** HIGH
- **Status:** PENDING
- **Deadline:** None

## Objective
Add a web UI to the existing Go CLI calculator so it can be accessed via a browser as an online tool.

## Background/Context
The calculator CLI (TASK-001) is complete and pushed to github.com/maxroyak/calculator. Now we need a web interface so it's accessible online. The Go binary should serve both CLI mode and web mode.

## Requirements
### Must Have (Required for completion)
- [ ] Follow all conventions in `shared/GOLANG_STANDARDS.md`
- [ ] Go HTTP server in `internal/server/` package
- [ ] HTML/JS frontend embedded in Go (use embed package or inline)
- [ ] CLI mode preserved — existing `calculator add 5 3` still works
- [ ] New `serve` subcommand starts the web server: `calculator serve`
- [ ] Port configurable via flag: `calculator serve --port 8080` (default 8080)
- [ ] Web UI supports all 4 operations (add, subtract, multiply, divide)
- [ ] Division by zero shows error in UI (no crash)
- [ ] Input validation on frontend and backend
- [ ] API endpoint: POST /api/calculate with JSON body {operation, a, b} → {result} or {error}
- [ ] Web page served at GET /
- [ ] Tests for HTTP handlers (table-driven, use httptest)
- [ ] Update README.md with web server usage
- [ ] Update Makefile if needed (add serve target)
- [ ] Append WORKLOG.md entries

### Should Have (Strongly desired)
- [ ] Clean, modern CSS styling (no frameworks — vanilla CSS)
- [ ] Responsive design (works on mobile)
- [ ] Keyboard support (Enter to calculate)

### Nice to Have (Optional)
- [ ] History of recent calculations in UI
- [ ] Health check endpoint GET /api/health

## Technical Specifications
### Files to Modify/Create
- `internal/server/server.go` — HTTP server, handlers, API endpoint
- `internal/server/server_test.go` — Handler tests using httptest
- `web/index.html` — Frontend HTML (or embed via go:embed)
- `web/style.css` — Styling
- `web/app.js` — Frontend logic
- `main.go` — Modify to add "serve" subcommand routing
- `README.md` — Add web server usage section
- `Makefile` — Add `serve` target (run server)
- `WORKLOG.md` — Append entries

### Architecture Notes
- Use Go 1.26+ embed package to bundle web/ files into the binary
- Reuse existing internal/calc package for all math operations
- API: POST /api/calculate → JSON {operation, a, b} → {result, error}
- Static files served from embedded FS
- Keep server logic in internal/server/, not in main.go
- main.go routes: if first arg is "serve" → start server, otherwise → CLI mode (existing behavior)

## Implementation Steps
1. Create internal/server/server.go with HTTP handlers
2. Create web/ directory with index.html, style.css, app.js
3. Use go:embed to bundle web files
4. Modify main.go to detect "serve" subcommand
5. Write handler tests in internal/server/server_test.go
6. Update README.md with serve instructions
7. Add `make serve` target to Makefile
8. Run go fmt, go vet, go test -race
9. Test server manually: start, curl API, check HTML response
10. Append WORKLOG.md entries

## Definition of Done
- [ ] Code compiles without errors
- [ ] All existing tests still pass
- [ ] New HTTP handler tests written and passing
- [ ] Code formatted (go fmt)
- [ ] Code passes vet (go vet)
- [ ] Follows GOLANG_STANDARDS.md
- [ ] CLI mode still works (backward compatible)
- [ ] Web UI works in browser
- [ ] Division by zero handled in both CLI and web
- [ ] README updated with web usage
- [ ] WORKLOG.md entries added
- [ ] Ready for QA review

## Testing Strategy
- httptest.NewRecorder for handler tests
- Test /api/calculate with all 4 operations
- Test division by zero returns error JSON
- Test invalid input returns error JSON
- Test / serves HTML content
- Test invalid operation returns error
- Table-driven tests for API endpoint

## Acceptance Criteria
1. `go build` succeeds
2. `go test -race ./...` passes (existing + new tests)
3. `go vet ./...` clean
4. `./build/calculator add 5 3` still outputs 8 (CLI backward compatible)
5. `./build/calculator serve` starts web server on :8080
6. `curl localhost:8080/` returns HTML page
7. `curl -X POST localhost:8080/api/calculate -d '{"operation":"add","a":5,"b":3}'` returns result
8. Division by zero returns error JSON, not crash
9. `make serve` works

## Questions/Blockers
| # | Question/Blocker | Status | Resolution |
|---|------------------|--------|------------|
| 1 | None | - | - |

## Progress Log
| Date | Time | Update |
|------|------|--------|
| 2026-08-15 | 23:55 | Task assigned by pm_bot |