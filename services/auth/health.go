package main

import (
	"encoding/json"
	"net/http"
)

// healthResponse describes the JSON body that /health returns.
// The `json:"..."` tags control the field names in the JSON output.
type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// healthHandler answers "is this service up?". A handler is any function that
// takes a ResponseWriter (to build the reply) and a Request (what was asked).
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(healthResponse{Status: "ok", Service: "auth"})
}
