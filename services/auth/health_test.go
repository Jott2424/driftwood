package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	// Build a fake request and a recorder that captures the reply.
	// No real server or network is involved.
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	healthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusOK)
	}

	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}

	var body healthResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("could not decode response body: %v", err)
	}
	if body.Status != "ok" {
		t.Errorf("status = %q, want %q", body.Status, "ok")
	}
	if body.Service != "auth" {
		t.Errorf("service = %q, want %q", body.Service, "auth")
	}
}

// TestHealthRoute checks the route as the server wires it up: GET works,
// POST is rejected. This shows the method matching in the router.
func TestHealthRoute(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)

	tests := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{"GET is allowed", http.MethodGet, http.StatusOK},
		{"POST is rejected", http.MethodPost, http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/health", nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status code = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
