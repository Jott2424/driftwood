package main

import (
	"context"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// The //go:embed line below makes the compiler copy every .sql file in the
// migrations folder into the program itself. The deployed .exe or binary then
// carries its own migrations, with no extra files to copy to a VM.
//
//go:embed migrations/*.sql
var migrationsFS embed.FS

// runMigrations applies any migrations that have not run yet. goose keeps a
// table called goose_db_version in the database to remember which ones have, so
// running this twice does nothing the second time.
func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	// goose speaks Go's standard database/sql interface, while we use pgx's pool.
	// stdlib.OpenDBFromPool wraps our pool so both can share it.
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
