package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// connectDB opens a pool of database connections and checks that the database
// answers. A pool keeps several connections open and hands them out to
// requests, so each request doesn't pay the cost of opening its own.
func connectDB(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		// %w wraps the original error so callers can still inspect it,
		// while we add context about what we were doing.
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	// pgxpool.New does not actually connect yet, so ping to find out now.
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}
