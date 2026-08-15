// Package server provides an HTTP server for the calculator web UI.
//
// It serves a static HTML/CSS/JS frontend (embedded via the web package)
// and a JSON API endpoint at POST /api/calculate that delegates math
// operations to the internal/calc package.
package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/maxroyak/calculator/internal/calc"
	"github.com/maxroyak/calculator/web"
)

// calculateRequest is the JSON body for POST /api/calculate.
type calculateRequest struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
}

// calculateResponse is the JSON response from POST /api/calculate.
// Exactly one of Result or Error is populated.
type calculateResponse struct {
	Result *float64 `json:"result,omitempty"`
	Error  string   `json:"error,omitempty"`
}

// Server wraps an http.Server configured with the calculator routes.
type Server struct {
	httpServer *http.Server
}

// New returns a Server that listens on the given address.
func New(addr string) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/calculate", handleCalculate)
	mux.HandleFunc("/api/health", handleHealth)
	mux.Handle("/", staticHandler())

	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

// ListenAndServe starts the HTTP server and blocks until it stops.
func (s *Server) ListenAndServe() error {
	log.Printf("calculator web server listening on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// staticHandler returns an http.Handler that serves embedded web/ files
// from the root path.
func staticHandler() http.Handler {
	return http.FileServer(http.FS(web.FS))
}

// handleHealth responds to GET /api/health with a simple JSON status.
func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

// handleCalculate processes POST /api/calculate requests.
func handleCalculate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(calculateResponse{Error: "method not allowed"})
		return
	}

	var req calculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(calculateResponse{Error: "invalid JSON body"})
		return
	}

	result, err := compute(req.Operation, req.A, req.B)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(calculateResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(calculateResponse{Result: &result})
}

// compute delegates to the calc package for the actual arithmetic.
func compute(operation string, a, b float64) (float64, error) {
	switch operation {
	case "add":
		return calc.Add(a, b), nil
	case "subtract":
		return calc.Subtract(a, b), nil
	case "multiply":
		return calc.Multiply(a, b), nil
	case "divide":
		return calc.Divide(a, b)
	default:
		return 0, fmt.Errorf("unknown operation %q: valid operations are add, subtract, multiply, divide", operation)
	}
}

// formatFloat converts a float64 to a clean string for display.
func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
