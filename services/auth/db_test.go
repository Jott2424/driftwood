package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// testDatabaseURL returns the address of a real Postgres for integration
// tests. If it isn't set, the test is skipped, so plain `go test` still works
// on a machine with no database running.
func testDatabaseURL(t *testing.T) string {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping integration test")
	}
	return url
}

func TestConnectDB(t *testing.T) {
	url := testDatabaseURL(t)

	pool, err := connectDB(context.Background(), url)
	if err != nil {
		t.Fatalf("connectDB: %v", err)
	}
	defer pool.Close()

	var one int
	if err := pool.QueryRow(context.Background(), "SELECT 1").Scan(&one); err != nil {
		t.Fatalf("query: %v", err)
	}
	if one != 1 {
		t.Errorf("SELECT 1 = %d, want 1", one)
	}
}

func TestConnectDBUnreachable(t *testing.T) {
	// Nothing listens on port 1, so connecting must fail, and quickly.
	_, err := connectDB(context.Background(), "postgres://nobody:nothing@127.0.0.1:1/none")
	if err == nil {
		t.Fatal("expected an error connecting to an unreachable database, got nil")
	}
}

func TestReadyWithRealDatabase(t *testing.T) {
	url := testDatabaseURL(t)

	pool, err := connectDB(context.Background(), url)
	if err != nil {
		t.Fatalf("connectDB: %v", err)
	}
	defer pool.Close()

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rec := httptest.NewRecorder()

	readyHandler(pool)(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusOK)
	}
}
