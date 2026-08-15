# Calculator

A simple command-line calculator written in Go that supports basic arithmetic operations: addition, subtraction, multiplication, and division.

## Features

- **Add** — Add two numbers
- **Subtract** — Subtract the second number from the first
- **Multiply** — Multiply two numbers
- **Divide** — Divide the first number by the second (with division-by-zero error handling)

## Project Structure

```
calculator/
├── main.go                    # CLI entry point and argument parsing
├── internal/
│   └── calc/
│       ├── calc.go            # Arithmetic operation functions
│       └── calc_test.go       # Table-driven unit tests
├── go.mod                     # Module definition
├── Makefile                   # Build automation
├── README.md                  # This file
└── WORKLOG.md                 # Append-only action log
```

## Prerequisites

- Go 1.26 or later

## Setup

Clone the repository:

```bash
git clone https://github.com/maxroyak/calculator.git
cd calculator
```

## Usage

### Build from source

```bash
make build
```

This produces a binary at `./build/calculator`.

### Run the calculator

```bash
./build/calculator <operation> <operand1> <operand2>
```

### Operations and Examples

```bash
./build/calculator add 5 3        # Output: 8
./build/calculator subtract 10 4  # Output: 6
./build/calculator multiply 6 7   # Output: 42
./build/calculator divide 10 2    # Output: 5
```

### Error Handling

```bash
./build/calculator divide 10 0
# Output: error: division by zero is not allowed

./build/calculator foo 1 2
# Output: error: unknown operation "foo": valid operations are add, subtract, multiply, divide

./build/calculator add abc 3
# Output: error: first operand: invalid number "abc"
```

### Show Help

```bash
./build/calculator --help
```

## Development

### Available Make Targets

| Target          | Description                         |
|-----------------|-------------------------------------|
| `make build`    | Build the application binary        |
| `make test`     | Run tests with verbose output       |
| `make race`     | Run tests with the race detector    |
| `make fmt`      | Format all Go source files          |
| `make vet`      | Run `go vet` on all packages        |
| `make lint`     | Run golangci-lint                   |
| `make clean`    | Remove build artifacts              |
| `make dev`      | Run fmt, vet, and test in sequence  |
| `make ci`       | Run fmt, vet, test, and race        |
| `make help`     | Show available targets              |

### Run Tests

```bash
make test          # standard
make race          # with race detector
go test -race ./...  # direct
```

### Format and Vet

```bash
make fmt
make vet
```

## Module

- **Module path:** `github.com/maxroyak/calculator`
- **Dependencies:** Go standard library only (no external dependencies)

## License

This is a test project for the openclaw-ai-dev-team workflow.