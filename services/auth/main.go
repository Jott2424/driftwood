package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Read the port from the environment, with a default for local use.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// The mux (router) maps a method and path to a handler.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second, // don't wait forever on slow clients
	}

	log.Printf("auth listening on :%s", port)
	// ListenAndServe blocks until the server stops. It always returns an error,
	// so if we get here, something went wrong.
	log.Fatal(server.ListenAndServe())
}
