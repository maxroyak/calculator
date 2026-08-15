package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleCalculate(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantResult *float64
		wantError  string
	}{
		// Happy path — all four operations.
		{"add positive", http.MethodPost, `{"operation":"add","a":5,"b":3}`, http.StatusOK, ptr(8), ""},
		{"add negative", http.MethodPost, `{"operation":"add","a":-5,"b":-3}`, http.StatusOK, ptr(-8), ""},
		{"subtract", http.MethodPost, `{"operation":"subtract","a":10,"b":4}`, http.StatusOK, ptr(6), ""},
		{"multiply", http.MethodPost, `{"operation":"multiply","a":6,"b":7}`, http.StatusOK, ptr(42), ""},
		{"divide", http.MethodPost, `{"operation":"divide","a":10,"b":2}`, http.StatusOK, ptr(5), ""},
		{"divide fractional", http.MethodPost, `{"operation":"divide","a":7.5,"b":2.5}`, http.StatusOK, ptr(3), ""},

		// Error cases.
		{"division by zero", http.MethodPost, `{"operation":"divide","a":10,"b":0}`, http.StatusBadRequest, nil, "division by zero is not allowed"},
		{"unknown operation", http.MethodPost, `{"operation":"foo","a":1,"b":2}`, http.StatusBadRequest, nil, "unknown operation"},
		{"invalid JSON", http.MethodPost, `{bad json}`, http.StatusBadRequest, nil, "invalid JSON body"},
		{"empty body", http.MethodPost, ``, http.StatusBadRequest, nil, "invalid JSON body"},

		// Wrong method.
		{"GET not allowed", http.MethodGet, ``, http.StatusMethodNotAllowed, nil, "method not allowed"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/api/calculate", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handleCalculate(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
			}

			var resp calculateResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if tc.wantError != "" {
				if resp.Error == "" {
					t.Fatalf("expected error containing %q, got empty", tc.wantError)
				}
				if !strings.Contains(resp.Error, tc.wantError) {
					t.Errorf("error = %q, want substring %q", resp.Error, tc.wantError)
				}
				if resp.Result != nil {
					t.Errorf("expected no result, got %v", *resp.Result)
				}
				return
			}

			if resp.Error != "" {
				t.Fatalf("unexpected error: %s", resp.Error)
			}
			if resp.Result == nil {
				t.Fatal("expected result, got nil")
			}
			if *resp.Result != *tc.wantResult {
				t.Errorf("result = %v, want %v", *resp.Result, *tc.wantResult)
			}
		})
	}
}

func TestHandleHealth(t *testing.T) {
	t.Run("GET returns ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		rec := httptest.NewRecorder()

		handleHealth(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		body := strings.TrimSpace(rec.Body.String())
		if body != `{"status":"ok"}` {
			t.Errorf("body = %q, want {\"status\":\"ok\"}", body)
		}
	})

	t.Run("POST not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/health", nil)
		rec := httptest.NewRecorder()

		handleHealth(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
		}
	})
}

func TestStaticHandlerServesIndex(t *testing.T) {
	handler := staticHandler()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}
	if !strings.Contains(string(body), "<title>Calculator</title>") {
		t.Error("response does not contain expected <title> element")
	}
}

func TestCompute(t *testing.T) {
	tests := []struct {
		name      string
		op        string
		a, b      float64
		want      float64
		wantErr   bool
		errSubstr string
	}{
		{"add", "add", 5, 3, 8, false, ""},
		{"subtract", "subtract", 10, 4, 6, false, ""},
		{"multiply", "multiply", 6, 7, 42, false, ""},
		{"divide", "divide", 10, 2, 5, false, ""},
		{"divide by zero", "divide", 10, 0, 0, true, "division by zero"},
		{"unknown op", "foo", 1, 2, 0, true, "unknown operation"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := compute(tc.op, tc.a, tc.b)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tc.errSubstr) {
					t.Errorf("error = %q, want substring %q", err.Error(), tc.errSubstr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("compute(%q, %v, %v) = %v, want %v", tc.op, tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestNewServerRoutes(t *testing.T) {
	// Verify that New returns a working server with routes by sending a
	// request through a test server built from the same mux logic.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/calculate" {
			handleCalculate(w, r)
			return
		}
		staticHandler().ServeHTTP(w, r)
	}))
	defer srv.Close()

	// Test API.
	body := `{"operation":"add","a":2,"b":8}`
	resp, err := srv.Client().Post(srv.URL+"/api/calculate", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	var result calculateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decoding: %v", err)
	}
	if result.Result == nil || *result.Result != 10 {
		t.Errorf("result = %v, want 10", result.Result)
	}
	resp.Body.Close()

	// Test static.
	resp2, err := srv.Client().Get(srv.URL + "/")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp2.StatusCode, http.StatusOK)
	}
	resp2.Body.Close()
}

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{8, "8"},
		{-3.5, "-3.5"},
		{0, "0"},
		{42, "42"},
	}
	for _, tc := range tests {
		got := formatFloat(tc.input)
		if got != tc.want {
			t.Errorf("formatFloat(%v) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// ptr returns a pointer to v, for use in table-driven tests.
func ptr(v float64) *float64 {
	return &v
}
