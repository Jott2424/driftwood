package main

import (
	"context"
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

	// The database address holds a password, so it has no default in the code.
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	pool, err := connectDB(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer pool.Close()

	// Bring the database structure up to date before serving any requests.
	if err := runMigrations(context.Background(), pool); err != nil {
		log.Fatalf("migrate database: %v", err)
	}

	// The mux (router) maps a method and path to a handler.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /ready", readyHandler(pool))

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second, // don't wait forever on slow clients
	}

	log.Printf("auth listening on :%s", port)
	// ListenAndServe blocks until the server stops. It always returns an error,
	// so if we get here, something went wrong. (log.Fatal skips the deferred
	// pool.Close, which is fine here because the process is exiting anyway.)
	log.Fatal(server.ListenAndServe())
}
