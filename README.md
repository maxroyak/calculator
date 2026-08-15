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
├── main.go                    # CLI entry point and serve subcommand routing
├── internal/
│   ├── calc/
│   │   ├── calc.go            # Arithmetic operation functions
│   │   └── calc_test.go       # Table-driven unit tests
│   └── server/
│       ├── server.go          # HTTP server, handlers, API endpoint
│       └── server_test.go     # Handler tests using httptest
├── web/
│   ├── embed.go               # go:embed directive for static assets
│   ├── index.html             # Frontend HTML
│   ├── style.css              # Styling (vanilla CSS)
│   └── app.js                 # Frontend logic (vanilla JS)
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

## Web Server

The calculator also includes a built-in web server with a browser-based UI.

### Start the Web Server

```bash
# Default port 8080
./build/calculator serve

# Custom port
./build/calculator serve --port 3000

# Or via Make
make serve
```

Then open `http://localhost:8080` in your browser.

### Web UI Features

- All four operations (add, subtract, multiply, divide)
- Clean, modern, responsive design (mobile-friendly)
- Keyboard support (Enter to calculate)
- Recent calculation history (last 10)
- Division by zero shows error message — no crash

### API Endpoint

The web server exposes a JSON API at `POST /api/calculate`:

```bash
curl -X POST localhost:8080/api/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"add","a":5,"b":3}'
# Response: {"result":8}
```

Division by zero returns an error:

```bash
curl -X POST localhost:8080/api/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"divide","a":10,"b":0}'
# Response: {"error":"division by zero is not allowed"}
```

### Health Check

```bash
curl localhost:8080/api/health
# Response: {"status":"ok"}
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
| `make serve`    | Build and start the web server (:8080) |
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