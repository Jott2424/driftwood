package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// pinger is anything that can check whether it is reachable. The real database
// pool does this, and so can a fake in tests. In Go, an interface is satisfied
// automatically by any type that has the right methods.
type pinger interface {
	Ping(ctx context.Context) error
}

type readyResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

// readyHandler answers "can this service do its job right now?". That differs
// from /health, which only says the process is alive. Orchestrators use the
// two differently: a failing health check means restart, a failing ready check
// means stop sending traffic for now.
//
// It returns a handler instead of being one, so it can carry the database it
// should check. This pattern is called a closure.
func readyHandler(db pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		w.Header().Set("Content-Type", "application/json")

		if err := db.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(readyResponse{Status: "unavailable", Database: "down"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(readyResponse{Status: "ok", Database: "up"})
	}
}
