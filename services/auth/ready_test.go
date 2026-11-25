package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// fakeDB satisfies the pinger interface without any real database.
type fakeDB struct {
	err error
}

func (f fakeDB) Ping(ctx context.Context) error { return f.err }

func TestReadyHandler(t *testing.T) {
	tests := []struct {
		name         string
		db           fakeDB
		wantStatus   int
		wantStatusJS string
		wantDBJS     string
	}{
		{"database up", fakeDB{}, http.StatusOK, "ok", "up"},
		{"database down", fakeDB{err: errors.New("connection refused")}, http.StatusServiceUnavailable, "unavailable", "down"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/ready", nil)
			rec := httptest.NewRecorder()

			readyHandler(tt.db)(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status code = %d, want %d", rec.Code, tt.wantStatus)
			}

			var body readyResponse
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatalf("could not decode response body: %v", err)
			}
			if body.Status != tt.wantStatusJS {
				t.Errorf("status = %q, want %q", body.Status, tt.wantStatusJS)
			}
			if body.Database != tt.wantDBJS {
				t.Errorf("database = %q, want %q", body.Database, tt.wantDBJS)
			}
		})
	}
}
