# Task Assignment: Go CLI Calculator

## Metadata
- **Task ID:** TASK-001
- **Project:** calculator
- **Assigned to:** dev_bot
- **Assigned by:** pm_bot
- **Date:** 2026-08-15
- **Priority:** MEDIUM
- **Status:** PENDING
- **Deadline:** None

## Objective
Build a simple CLI calculator in Go that supports add, subtract, multiply, and divide operations.

## Background/Context
This is a test project to verify the team workflow and push a working Go project to GitHub (github.com/maxroyak/calculator).

## Requirements
### Must Have (Required for completion)
- [ ] Follow all conventions in `shared/GOLANG_STANDARDS.md`
- [ ] Project structure matches `shared/GOLANG_PROJECT_TEMPLATE.md`
- [ ] Include Makefile based on `shared/MAKEFILE_TEMPLATE`
- [ ] Include README.md with setup and usage
- [ ] CLI calculator supporting: add, subtract, multiply, divide
- [ ] Handle division by zero gracefully
- [ ] Table-driven unit tests for all operations
- [ ] Input validation with meaningful error messages

### Should Have (Strongly desired)
- [ ] Support for chained operations (e.g., 5 + 3 * 2)
- [ ] Interactive mode (REPL)

### Nice to Have (Optional)
- [ ] Support for additional operations (power, sqrt, modulo)

## Technical Specifications
### Files to Modify/Create
- `go.mod` - Module definition (github.com/maxroyak/calculator)
- `main.go` - Entry point, CLI argument parsing
- `internal/calc/calc.go` - Calculator logic (operations)
- `internal/calc/calc_test.go` - Table-driven tests
- `Makefile` - Build automation (from MAKEFILE_TEMPLATE)
- `README.md` - Project documentation
- `WORKLOG.md` - Already created by pm_bot

### Dependencies
- Go standard library only (no external deps)

### Architecture Notes
- Keep main package minimal, delegate logic to internal/calc package
- Use flag package for CLI argument parsing
- Operations as exported functions: Add, Subtract, Multiply, Divide

## Implementation Steps
1. Initialize go module: `go mod init github.com/maxroyak/calculator`
2. Create internal/calc package with operation functions
3. Create main.go with CLI argument parsing (operation + two operands)
4. Write table-driven tests in internal/calc/calc_test.go
5. Create Makefile from template
6. Write README.md with usage examples
7. Run go fmt, go vet, go test -race
8. Append WORKLOG.md entries

## Definition of Done
- [ ] Code compiles without errors
- [ ] Unit tests written and passing
- [ ] Code formatted (`go fmt`)
- [ ] Code passes vet (`go vet`)
- [ ] Follows GOLANG_STANDARDS.md
- [ ] No critical TODO comments left
- [ ] Documentation updated
- [ ] **WORKLOG.md entry added for task completion**
- [ ] Ready for QA review

## Testing Strategy
- Table-driven tests for each operation (positive cases)
- Division by zero test
- Invalid input handling test
- Edge cases: negative numbers, floating point precision

## Acceptance Criteria
1. `go build` succeeds
2. `go test -race ./...` passes
3. `go vet ./...` clean
4. `make build` and `make test` work
5. Division by zero returns error, not panic
6. README has clear usage instructions

## Questions/Blockers
| # | Question/Blocker | Status | Resolution |
|---|------------------|--------|------------|
| 1 | None | - | - |

## Progress Log
| Date | Time | Update |
|------|------|--------|
| 2026-08-15 | 00:00 | Task assigned by pm_bot |