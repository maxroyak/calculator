# QA Report: Go CLI Calculator

## Metadata
- **Report ID:** QA-001
- **From:** qa_bot
- **To:** pm_bot
- **Project:** calculator
- **Date:** 2026-08-15
- **Review Duration:** ~0.2 hours
- **Handover Reference:** WORKLOG 2026-08-15 23:49 (dev_bot TASK_COMPLETE)
- **Status:** APPROVED_WITH_COMMENTS

## Review Summary
The Go CLI calculator is a well-structured, minimal project that follows the team's Go standards closely. All toolchain checks pass: `go fmt`, `go vet`, `go test -race`, and `go build` are clean. The `internal/calc` package achieves 100% test coverage with table-driven tests covering positive, negative, fractional, zero, and error cases. Division by zero correctly returns an error rather than panicking. No security vulnerabilities were identified. The only notable gap is the absence of unit tests for the `main` package (functions `run`, `calculate`, `parseFloat`, `printUsage`), which reports 0% coverage.

### Overall Assessment
| Aspect | Rating | Notes |
|--------|--------|-------|
| Code Quality | Excellent | Clean, idiomatic Go; proper error wrapping; early returns |
| Test Coverage | Good | 100% on `internal/calc`; 0% on `main` package (no test files) |
| Security | Excellent | Input validation present; no external deps; no injection surface |
| Performance | Excellent | No concerns for a CLI calculator; no unnecessary allocations |
| Documentation | Good | Godoc comments on all exported symbols; README thorough; missing DEV_HANDOVER.md |

## Issues Found

### Critical Issues (MUST FIX - Blocks approval)
| # | Issue | File | Line | Severity | Description | Recommendation |
|---|-------|------|------|----------|-------------|----------------|
| — | — | — | — | — | No critical issues found | — |

**Critical Issues Count:** 0
**Status:** N/A

### Important Issues (SHOULD FIX - Strongly recommended)
| # | Issue | File | Line | Severity | Description | Recommendation |
|---|-------|------|------|----------|-------------|----------------|
| 1 | QA-IMP-001 | main.go | — | Important | No test file for the `main` package; `run`, `calculate`, `parseFloat`, and `printUsage` have 0% coverage. Standards require >80% on critical paths and `go test ./...` reports `?[no test files]` for the root package. | Add `main_test.go` with table-driven tests for `calculate` (all operations + unknown op), `parseFloat` (valid/invalid), `run` (no args, help, valid op, insufficient args, invalid operand), and `printUsage`. |
| 2 | QA-IMP-002 | go.mod | 3 | Important | `go.mod` declares `go 1.25.0` but README states "Go 1.26 or later" and the toolchain in use is `go1.26.0`. The version mismatch is confusing and could cause CI failures on environments with only Go 1.25. | Align: either update `go.mod` to `go 1.26.0` to match the README, or update the README prerequisite to "Go 1.25 or later" to match `go.mod`. |
| 3 | QA-IMP-003 | — | — | Important | No `DEV_HANDOVER.md` exists in the project directory. The HANDOVER_PROTOCOL.md (§3) requires dev_bot to create this document before QA review. | Request dev_bot to create `DEV_HANDOVER.md` per the template in HANDOVER_PROTOCOL.md §3.3. |

**Important Issues Count:** 3
**Status:** 3 remaining

### Minor Issues (NICE TO FIX - Non-blocking)
| # | Issue | File | Line | Severity | Description | Recommendation |
|---|-------|------|------|----------|-------------|----------------|
| 1 | QA-MIN-001 | Makefile | 6-7 | Minor | `GO_FILES` and `PACKAGES` variables are defined but never used by any target. | Remove unused variables or use them in `lint`/`fmt` targets to be consistent with the template's intent. |
| 2 | QA-MIN-002 | Makefile | 16 | Minor | Build target builds from root (`.`) rather than `./cmd/$(BINARY_NAME)` as in the template. This is intentional for this simple project (main.go at root, per template's "for simple apps" allowance) but deviates from the template's default. | Document this as a conscious decision in a comment, or consider moving `main.go` to `cmd/calculator/` for full template alignment. |
| 3 | QA-MIN-003 | main.go | 77-95 | Minor | `printUsage` returns `errors.New("no operation provided")` — it both prints usage to stdout and returns an error. This dual behavior means calling `--help` prints usage but also exits with code 1 and an error message, which is unusual for `--help`/`-h` flags. | Consider returning `nil` for explicit `--help`/`-h` invocations (exit 0) and only returning the error for the no-args case. |
| 4 | QA-MIN-004 | README.md | 29 | Minor | README states "Go 1.26 or later" but `go.mod` says `go 1.25.0`. See QA-IMP-002. | Fix in conjunction with QA-IMP-002. |
| 5 | QA-MIN-005 | WORKLOG.md | 6 | Minor | The entry on line 6 is missing the opening `---` separator: the line reads `---2026-08-15 23:46` (merged with the previous entry's closing `---`). | Add a newline and `---` separator to properly delimit the entry per the WORKLOG format. |

**Minor Issues Count:** 5

## Security Findings

### Vulnerabilities
| # | Vulnerability | CVSS | Description | Mitigation |
|---|---------------|------|-------------|------------|
| — | — | — | No vulnerabilities identified | — |

### Security Concerns
- **Input validation:** All numeric inputs are validated via `parseFloat` with clear error messages on failure. ✓
- **Division by zero:** Returns a descriptive error, no panic. ✓
- **Unknown operation:** Returns a descriptive error listing valid operations. ✓
- **No external dependencies:** The project uses only the Go standard library, eliminating supply-chain risk. ✓
- **No injection surface:** No eval, no shell execution, no file I/O, no network calls. ✓
- **Error messages:** Do not leak sensitive internal state. ✓
- **golangci-lint / gosec / govulncheck:** Not available in the environment. Recommend running these in CI when the tools are installed. The Makefile includes a `lint` target referencing `golangci-lint`.

## Test Quality Assessment

### Coverage Analysis
- **Current coverage:** `internal/calc` = 100%; `main` package = 0% (no test files)
- **Target coverage:** >80% on critical paths (per GOLANG_STANDARDS.md)
- **Gap:** The `main` package's CLI logic (`run`, `calculate`, `parseFloat`, `printUsage`) is untested. While the `internal/calc` library is fully covered, the argument parsing and operation dispatch logic in `main.go` — which is user-facing — has no automated tests.

### Test Quality
| Aspect | Rating | Notes |
|--------|--------|-------|
| Test independence | Good | Each test case is self-contained; no shared mutable state |
| Edge case coverage | Good | Covers positive, negative, mixed signs, zero, fractional, division-by-zero error |
| Test maintainability | Good | Table-driven with named subtests; clear error messages |

### Missing Tests
| Test | File | Priority | Reason Missing |
|------|------|----------|----------------|
| TestCalculate | main_test.go | High | No test file for `main` package; `calculate` dispatch is critical user-facing logic |
| TestParseFloat | main_test.go | High | Input validation for numeric parsing; no test file exists |
| TestRun | main_test.go | Medium | Integration of argument parsing + operation dispatch; no test file exists |
| TestPrintUsage | main_test.go | Low | Usage output; no test file exists |

## Performance Findings

### Observations
No performance concerns. The calculator is a single-shot CLI with no loops, no goroutines, no I/O beyond stdin/stdout. All operations are O(1) arithmetic.

### Recommendations
None.

## Code Review Comments

### General Comments
The code is clean, idiomatic Go that follows Effective Go conventions. Error handling follows the standards: errors are checked explicitly, wrapped with context using `fmt.Errorf` with `%w`, and returned early. The `internal/calc` package is well-documented with godoc comments on the package and all exported functions. The `main` package has a package-level godoc comment. Import grouping is correct (stdlib, then local). Line lengths are well within 120 characters.

### Specific Comments
| # | File | Line | Comment |
|---|------|------|---------|
| 1 | main.go | 1-2 | Package comment is good — describes purpose and operations. |
| 2 | main.go | 36 | Good error wrapping with `%w` for operand parse errors. |
| 3 | main.go | 46 | `calculate` error is returned without wrapping context. Consider `fmt.Errorf("calculate: %w", err)` for consistency, though the underlying errors are already descriptive. |
| 4 | calc.go | 1-5 | Excellent package-level godoc documenting the division-by-zero contract. |
| 5 | calc.go | 26-29 | `Divide` correctly returns error on zero divisor — no panic. Clean implementation. |
| 6 | calc_test.go | 1 | No package comment. Since this is an internal test file (`package calc`), this is acceptable but a brief comment would be consistent with the standards. |

## Recommendations

### Required Actions
(None — no critical issues block approval.)

### Suggested Improvements
1. **Add `main_test.go`** (QA-IMP-001): Create table-driven tests for `calculate`, `parseFloat`, and `run` to bring the `main` package to >80% coverage. This is the most impactful improvement.
2. **Fix go.mod / README version mismatch** (QA-IMP-002): Align the Go version in both files.
3. **Create DEV_HANDOVER.md** (QA-IMP-003): Per the handover protocol, dev_bot should provide this before QA review.
4. **Improve `--help` behavior** (QA-MIN-003): Return `nil` for explicit help flags so exit code is 0.
5. **Fix WORKLOG.md formatting** (QA-MIN-005): Add missing `---` separator on line 6.
6. **Remove unused Makefile variables** (QA-MIN-001): Clean up `GO_FILES` and `PACKAGES` if not used.
7. **Install and run security tools in CI**: `golangci-lint`, `gosec`, and `govulncheck` are not available in the current environment. The Makefile has a `lint` target; recommend adding `gosec` and `govulncheck` targets and running them in CI.

### Best Practices Not Followed
| Practice | Location | Recommendation |
|----------|----------|----------------|
| >80% coverage on critical paths | main.go (entire file) | Add tests for `run`, `calculate`, `parseFloat` |
| DEV_HANDOVER.md before QA | project root | dev_bot should create handover document |
| WORKLOG.md entry format | WORKLOG.md:6 | Fix missing `---` delimiter |

## Approval Status

### Decision
**Status:** APPROVED_WITH_COMMENTS

The project is functionally correct, secure, and well-structured. All toolchain checks pass. The `internal/calc` package has excellent test coverage. The main gaps are: (1) no tests for the `main` package's CLI logic, (2) a go.mod/README version mismatch, and (3) a missing DEV_HANDOVER.md. None of these are blocking — the calculator works correctly and safely — but they should be addressed for full standards compliance.

### Conditions
1. Add `main_test.go` with tests for `calculate`, `parseFloat`, and `run` to achieve >80% coverage on the `main` package (target: before next release).
2. Align Go version between `go.mod` and `README.md`.
3. Create `DEV_HANDOVER.md` per HANDOVER_PROTOCOL.md §3.3.

### Next Steps
- [ ] dev_bot addresses QA-IMP-001 (add main package tests)
- [ ] dev_bot addresses QA-IMP-002 (fix go version mismatch)
- [ ] dev_bot addresses QA-IMP-003 (create DEV_HANDOVER.md)
- [ ] qa_bot re-reviews after fixes (if needed)
- [ ] Minor issues can be addressed opportunistically

## Time Spent
| Activity | Time |
|----------|------|
| Code review | 0.1 hours |
| Security review | 0.05 hours |
| Testing (running toolchain checks) | 0.02 hours |
| Documentation (writing QA_REPORT.md) | 0.05 hours |
| **Total** | **~0.2 hours** |

## Sign-Off

**QA Engineer:** qa_bot
**Date:** 2026-08-15
**Decision:** APPROVED_WITH_COMMENTS

**WORKLOG:** Reviewer confirms WORKLOG.md entries created for review start and completion.

---

## Verification Results Summary

| Check | Command | Result |
|-------|---------|--------|
| Format | `go fmt ./...` | ✅ Clean (no files needed formatting) |
| Vet | `go vet ./...` | ✅ Clean (exit 0) |
| Test (race) | `go test -race ./...` | ✅ All pass (21 subtests across 4 test functions) |
| Coverage | `go test -cover ./...` | `internal/calc`: 100% · `main`: 0% (no test files) |
| Build | `go build -o /dev/null .` | ✅ Success (exit 0) |
| Binary — add | `calculator add 5 3` | ✅ Output: `8` |
| Binary — subtract | `calculator subtract 10 4` | ✅ Output: `6` |
| Binary — multiply | `calculator multiply 6 7` | ✅ Output: `42` |
| Binary — divide | `calculator divide 10 2` | ✅ Output: `5` |
| Binary — div by zero | `calculator divide 10 0` | ✅ Error: `division by zero is not allowed` (exit 1) |
| Binary — unknown op | `calculator foo 1 2` | ✅ Error: `unknown operation "foo"...` (exit 1) |
| Binary — invalid operand | `calculator add abc 3` | ✅ Error: `first operand: invalid number "abc"` (exit 1) |
| Binary — insufficient args | `calculator add 5` | ✅ Error: `operation "add" requires two numeric operands` (exit 1) |
| Binary — no args / help | `calculator` / `calculator --help` | ✅ Prints usage (exit 1) |
| golangci-lint | N/A | ⚠️ Tool not installed in environment |
| gosec | N/A | ⚠️ Tool not installed in environment |
| govulncheck | N/A | ⚠️ Tool not installed in environment |

## Standards Compliance Checklist

| Standard (from GOLANG_STANDARDS.md) | Status | Notes |
|--------------------------------------|--------|-------|
| Follow Effective Go | ✅ | Idiomatic Go throughout |
| Use gofmt | ✅ | `go fmt ./...` clean |
| Group imports (stdlib, third-party, local) | ✅ | stdlib block, then local import |
| Lines < 120 chars | ✅ | All lines within limit |
| MixedCaps for exported names | ✅ | `Add`, `Subtract`, `Multiply`, `Divide` |
| camelCase for unexported names | ✅ | `run`, `calculate`, `parseFloat`, `printUsage` |
| Check all errors explicitly | ✅ | All errors checked |
| Return early on errors | ✅ | Early returns throughout `run` |
| Wrap errors with context | ✅ | `fmt.Errorf` with `%w` |
| Don't ignore errors with `_` | ✅ | No ignored errors |
| Meaningful error messages | ✅ | Descriptive messages |
| Table-driven tests | ✅ | All 4 test functions are table-driven |
| Test positive + negative cases | ✅ | Both covered |
| >80% coverage on critical paths | ⚠️ | 100% on `calc`, 0% on `main` package |
| Validate all inputs | ✅ | `parseFloat` validates numeric input |
| Godoc comments on exported symbols | ✅ | All exported functions documented |
| README with setup and usage | ✅ | Comprehensive README |
| WORKLOG.md maintained | ✅ | Entries present (minor formatting issue) |
| Makefile present | ✅ | Follows MAKEFILE_TEMPLATE with appropriate adjustments |