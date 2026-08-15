# QA Report: Web UI Additions (TASK-002)

## Metadata
- **Report ID:** QA-002
- **From:** qa_bot
- **To:** pm_bot
- **Project:** calculator
- **Date:** 2026-08-15
- **Review Duration:** 0.5 hours
- **Task Reference:** TASK-002 (Web UI additions)
- **Status:** APPROVED_WITH_COMMENTS

## Review Summary
The TASK-002 web UI additions are well-structured, follow Go conventions closely, and all toolchain checks pass (go fmt, go vet, go test -race, go build). The HTTP server correctly implements the JSON API with proper input validation, division-by-zero returns an error JSON (no crash), and CLI backward compatibility is fully preserved. Server package test coverage is 81.6% with table-driven tests covering positive, negative, and edge cases. Two important issues (missing main package tests, no request body size limit) and four minor issues are noted below.

### Overall Assessment
| Aspect | Rating | Notes |
|--------|--------|-------|
| Code Quality | Good | Clean, idiomatic Go. Godoc comments present. Proper error handling. Unused function `formatFloat` is dead code. |
| Test Coverage | Good | internal/server 81.6%, internal/calc 100%. main package 0%. Tests are table-driven with positive + negative + edge cases. |
| Security | Good | Input validation present. No injection surface (JSON API, static files only). Missing request body size limit (DoS vector). No CORS headers (acceptable for same-origin). |
| Performance | Good | Embed.FS is efficient. No allocations concerns. http.FileServer for static serving is standard. |
| Documentation | Good | README updated with web server section. Godoc on all exported types/functions. Package comments present. |

## Issues Found

### Critical Issues (MUST FIX - Blocks approval)
| # | Issue | File | Line | Severity | Description | Recommendation |
|---|-------|------|------|----------|-------------|----------------|
| (none) | | | | | | |

**Critical Issues Count:** 0
**Status:** N/A

### Important Issues (SHOULD FIX - Strongly recommended)
| # | Issue | File | Line | Severity | Description | Recommendation |
|---|-------|------|------|----------|-------------|----------------|
| 1 | QA-002-I1 | internal/server/server.go | 87 | Important | No request body size limit on JSON decode. `json.NewDecoder(r.Body).Decode(&req)` reads an unbounded request body. A malicious client could send a very large body to exhaust server memory (DoS). | Wrap `r.Body` with `io.LimitReader(r.Body, maxBodySize)` before decoding, e.g. `r.Body = http.MaxBytesReader(w, r.Body, 1<<20)` (1 MB). Return 413 or 400 if exceeded. |
| 2 | QA-002-I2 | main.go | 1 | Important | main package has 0% test coverage. Functions `run`, `runCLI`, `runServe`, `calculate`, `parseFloat`, `printUsage` are untested. This was flagged in QA-001 and remains unaddressed. The `runServe` function in particular has no test. | Add `main_test.go` with table-driven tests for `calculate`, `parseFloat`, and `run`/`runCLI` (the latter can be tested by calling `run` with arg slices and checking output/error). |

**Important Issues Count:** 2
**Status:** 0 fixed, 2 remaining

### Minor Issues (NICE TO FIX - Non-blocking)
| # | Issue | File | Line | Severity | Description | Recommendation |
|---|-------|------|------|----------|-------------|----------------|
| 1 | QA-002-M1 | internal/server/server.go | 121-123 | Minor | `formatFloat` function is defined but never called anywhere in the codebase (dead code). It is tested but serves no runtime purpose. | Either use it in the response (e.g., return a formatted string alongside the numeric result) or remove it and its test. |
| 2 | QA-002-M2 | internal/server/server_test.go | 170-179 | Minor | `TestNewServerRoutes` does not actually test `server.New()`. It constructs a custom `httptest.Server` with inline handler logic that duplicates the mux wiring, so `New()` and its `ListenAndServe()` are never exercised. | Consider testing `New()` by extracting the `http.Handler` (mux) from the Server struct (e.g., via a `Handler()` accessor) and passing it to `httptest.NewServer`. |
| 3 | QA-002-M3 | web/app.js | 46 | Minor | `addHistory` uses `li.innerHTML` to construct history entries. While the inputs are user-typed numbers (not strings from the user), using `innerHTML` with string concatenation is an XSS-prone pattern. The `result` value comes from the server JSON response and `a`/`b` from input fields — these are numeric, so exploitation is unlikely, but the pattern is unsafe by default. | Use `textContent` for the entry text and `createElement`/`appendChild` for the result span, or sanitize values before inserting via `innerHTML`. |
| 4 | QA-002-M4 | internal/server/server.go | 95-96 | Minor | Compute errors (e.g., unknown operation, division by zero) all return HTTP 400 Bad Request. While this is defensible, a case could be made for returning 422 Unprocessable Entity for unknown operations (the request is well-formed JSON but semantically invalid). This is a style preference, not a bug. | Consider differentiating 400 (malformed request) from 422 (semantically invalid) if the API evolves. Low priority. |

**Minor Issues Count:** 4

## Security Findings

### Vulnerabilities
| # | Vulnerability | CVSS | Description | Mitigation |
|---|---------------|------|-------------|------------|
| 1 | Unbounded Request Body | 3.7 (Low) | `handleCalculate` reads request body without a size limit. A client can send arbitrarily large JSON to consume server memory. | Use `http.MaxBytesReader` or `io.LimitReader` to cap body size (e.g., 1 MB). |

### Security Concerns
- **No CORS headers:** The API does not set CORS headers. This is acceptable since the web UI is same-origin (served from the same server). If the API is intended for cross-origin use, CORS would need to be configured.
- **No rate limiting:** The server has no rate limiting. For a local calculator tool this is acceptable; for production deployment it would be needed.
- **No TLS:** The server is plain HTTP. Acceptable for local development; document that a reverse proxy should handle TLS for any non-local deployment.
- **Static file serving:** `http.FileServer` with `http.FS(web.FS)` serves embedded files only — no path traversal risk since `embed.FS` does not allow escaping the embedded directory.
- **Input validation:** JSON request is properly decoded into a typed struct. Operation is validated via switch/default. Numbers are `float64` (Go handles invalid JSON decode with error). No SQL, shell, or template injection surface.
- **XSS in frontend:** `innerHTML` usage in `app.js:46` is a minor concern (see QA-002-M3). Values are numeric so exploitation requires a modified server response.

## Test Quality Assessment

### Coverage Analysis
- **Current coverage:** internal/server 81.6%, internal/calc 100.0%, main 0.0%, web N/A
- **Target coverage:** >80% on critical paths (per GOLANG_STANDARDS.md)
- **Gap:** main package at 0% (important), `New()` and `ListenAndServe()` at 0% in server package

### Test Quality
| Aspect | Rating | Notes |
|--------|--------|-------|
| Test independence | Good | Tests use httptest.NewRecorder and httptest.NewRequest — no shared state, no external dependencies. |
| Edge case coverage | Good | Division by zero, invalid JSON, empty body, wrong HTTP method, unknown operation all covered. |
| Test maintainability | Good | Table-driven tests with clear names. `ptr()` helper for float64 pointers. Subtests via `t.Run`. |

### Missing Tests
| Test | File | Priority | Reason Missing |
|------|------|----------|----------------|
| `server.New()` route wiring | server_test.go | Medium | `TestNewServerRoutes` tests a hand-built mux, not the actual `New()` function |
| `server.ListenAndServe()` | server_test.go | Low | Hard to test without port binding; acceptable to skip |
| main package functions | main_test.go | High | No test file exists for main package (flagged in QA-001) |
| Request body size limit | server_test.go | Medium | No limit exists to test (see QA-002-I1) |
| Very large number handling | server_test.go | Low | No test for float64 overflow/inf/NaN edge cases |

## Performance Findings

### Observations
- `embed.FS` is zero-copy at runtime — static assets are embedded in the binary, no disk I/O.
- `http.FileServer` is the standard Go approach for static serving with proper caching headers.
- JSON encoding/decoding uses `json.NewEncoder`/`json.NewDecoder` (streaming) — efficient for small payloads.
- No goroutines, no mutexes — the server is purely request-scoped, no race conditions possible.

### Recommendations
- No performance optimizations needed for this use case.

## Code Review Comments

### General Comments
- Code is clean, idiomatic, and well-organized. Package structure follows the project template.
- Error handling is consistent: errors are checked, wrapped with context where needed, and returned to the caller.
- Import grouping is correct: standard library first, then local packages.
- Godoc comments are present on all exported symbols (`Server`, `New`, `ListenAndServe`) and on both packages (`server`, `web`).
- The `compute` function in `server.go` duplicates the `calculate` function in `main.go`. Both are switch statements over the same operations. This is a minor DRY violation but acceptable given the different return types and package boundaries.

### Specific Comments
| # | File | Line | Comment |
|---|------|------|---------|
| 1 | server.go | 62-63 | `staticHandler()` returns a new `http.FileServer` on every call. In `New()` it's called once, so this is fine. But `TestNewServerRoutes` calls it per-request in the test harness. Minor inefficiency in tests only. |
| 2 | server.go | 87 | `json.NewDecoder(r.Body)` does not limit body size. See QA-002-I1. |
| 3 | server.go | 121-123 | `formatFloat` is dead code. See QA-002-M1. |
| 4 | app.js | 87 | `Number(data.result).toString()` is used to format the result client-side. The server already returns a float64; `JSON.parse` handles it. The `Number()` call is redundant but harmless. |
| 5 | main.go | 69-79 | `runServe` does not handle the case where `server.New(addr)` returns nil (it can't with current code, but the pattern is fragile). Consider adding a nil check or making `New` return an error. |
| 6 | main.go | 34-36 | `op == "serve"` is checked with string equality. This is fine for a simple CLI. If the CLI grows, consider using a subcommand framework (e.g., `flag.NewFlagSet` per subcommand). |

## Recommendations

### Required Actions
(No required actions — no critical issues found.)

### Suggested Improvements
1. Add `http.MaxBytesReader` to `handleCalculate` to limit request body size (QA-002-I1).
2. Add `main_test.go` with tests for `calculate`, `parseFloat`, and `run`/`runCLI` (QA-002-I2).
3. Remove or use the dead `formatFloat` function (QA-002-M1).
4. Fix `TestNewServerRoutes` to actually test `server.New()` (QA-002-M2).
5. Replace `innerHTML` with `textContent`/`createElement` in `app.js` (QA-002-M3).

### Best Practices Not Followed
| Practice | Location | Recommendation |
|----------|----------|----------------|
| Request body size limiting | server.go:87 | Use `http.MaxBytesReader` before decoding JSON |
| Dead code removal | server.go:121-123 | Remove unused `formatFloat` or integrate it |
| DRY (calculator logic) | main.go:81 / server.go:105 | Consider a shared `calc.Operate(op, a, b)` function to eliminate the duplicated switch |

## Approval Status

### Decision
**Status:** APPROVED_WITH_COMMENTS

### Conditions
1. Address QA-002-I1 (request body size limit) before any non-local deployment.
2. Address QA-002-I2 (main package tests) in a future iteration — carries over from QA-001.

### Next Steps
- [ ] dev_bot to address QA-002-I1 (add MaxBytesReader) — recommended before next release
- [ ] dev_bot to address QA-002-I2 (add main_test.go) — carried over from QA-001
- [ ] pm_bot to decide on minor issues (M1–M4) — defer or schedule
- [ ] Consider extracting shared `calc.Operate()` to eliminate duplicate switch in main.go and server.go

## Toolchain Verification Results

| Check | Command | Result |
|-------|---------|--------|
| Format | `go fmt ./...` | PASS (exit 0, no changes) |
| Vet | `go vet ./...` | PASS (exit 0, clean) |
| Test (race) | `go test -race ./...` | PASS (calc 1.011s, server 1.025s) |
| Build | `go build -o build/calculator .` | PASS (exit 0) |
| Coverage | `go test -cover ./...` | calc 100%, server 81.6%, main 0% |

## Runtime Verification Results

| Test | Result |
|------|--------|
| CLI: add 5 3 | ✅ Output: 8, exit 0 |
| CLI: subtract 10 4 | ✅ Output: 6, exit 0 |
| CLI: multiply 6 7 | ✅ Output: 42, exit 0 |
| CLI: divide 10 2 | ✅ Output: 5, exit 0 |
| CLI: divide 10 0 | ✅ Error: "division by zero is not allowed", exit 1 |
| CLI: foo 1 2 | ✅ Error: "unknown operation", exit 1 |
| CLI: no args | ✅ Usage printed, exit 1 |
| API: POST /api/calculate add | ✅ `{"result":8}` |
| API: POST /api/calculate divide by zero | ✅ `{"error":"division by zero is not allowed"}` |
| API: POST /api/calculate unknown op | ✅ `{"error":"unknown operation ..."}` |
| API: POST /api/calculate invalid JSON | ✅ `{"error":"invalid JSON body"}` |
| API: POST /api/calculate empty body | ✅ `{"error":"invalid JSON body"}` |
| API: POST /api/calculate missing fields | ✅ `{"result":0}` (zero-value for missing a/b) |
| API: GET /api/calculate (wrong method) | ✅ HTTP 405 |
| API: GET /api/health | ✅ `{"status":"ok"}` |
| API: POST /api/health (wrong method) | ✅ HTTP 405 |
| Static: GET / | ✅ HTTP 200 (index.html) |
| Static: GET /style.css | ✅ HTTP 200 |
| Static: GET /app.js | ✅ HTTP 200 |
| Static: GET /nonexistent | ✅ HTTP 404 |

## Time Spent
| Activity | Time |
|----------|------|
| Code review | 0.2 hours |
| Security review | 0.1 hours |
| Testing (toolchain + runtime) | 0.1 hours |
| Documentation | 0.1 hours |
| **Total** | 0.5 hours |

## Sign-Off

**QA Engineer:** qa_bot
**Date:** 2026-08-15
**Decision:** APPROVED_WITH_COMMENTS

**WORKLOG:** Reviewer confirms WORKLOG.md entry created with findings summary.