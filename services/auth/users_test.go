package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// testPool connects to the test database and makes sure migrations have run.
func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := testDatabaseURL(t)

	pool, err := connectDB(context.Background(), url)
	if err != nil {
		t.Fatalf("connectDB: %v", err)
	}
	t.Cleanup(pool.Close)

	if err := runMigrations(context.Background(), pool); err != nil {
		t.Fatalf("runMigrations: %v", err)
	}
	return pool
}

// uniqueName makes names that don't collide between test runs.
func uniqueName(prefix string) string {
	return fmt.Sprintf("%s%d", prefix, time.Now().UnixNano())
}

// createTestUser inserts a user and deletes it when the test ends.
func createTestUser(t *testing.T, repo *UserRepo, pool *pgxpool.Pool, email, username string) User {
	t.Helper()
	u, err := repo.Create(context.Background(), email, username, "hash-not-real")
	if err != nil {
		t.Fatalf("Create(%q, %q): %v", email, username, err)
	}
	t.Cleanup(func() {
		pool.Exec(context.Background(), "DELETE FROM users WHERE id = $1", u.ID)
	})
	return u
}

func TestMigrationsCreateUsersTable(t *testing.T) {
	pool := testPool(t)

	rows, err := pool.Query(context.Background(),
		`SELECT column_name FROM information_schema.columns
		 WHERE table_name = 'users' ORDER BY column_name`)
	if err != nil {
		t.Fatalf("query columns: %v", err)
	}
	defer rows.Close()

	var got []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("scan: %v", err)
		}
		got = append(got, name)
	}

	want := "created_at,email,id,password_hash,username"
	if strings.Join(got, ",") != want {
		t.Errorf("users columns = %v, want %s", got, want)
	}
}

func TestRunMigrationsTwice(t *testing.T) {
	pool := testPool(t)

	// testPool already ran them once. A second run must succeed and change nothing.
	if err := runMigrations(context.Background(), pool); err != nil {
		t.Fatalf("second runMigrations: %v", err)
	}
}

func TestUserRepoCreateAndGet(t *testing.T) {
	pool := testPool(t)
	repo := NewUserRepo(pool)

	email := uniqueName("ann") + "@example.com"
	username := uniqueName("ann")
	created := createTestUser(t, repo, pool, email, username)

	if created.ID == "" {
		t.Error("Create returned an empty ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("Create returned a zero CreatedAt")
	}

	got, err := repo.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Email != email || got.Username != username {
		t.Errorf("GetByID = %+v, want email %q username %q", got, email, username)
	}
}

func TestUserRepoDuplicates(t *testing.T) {
	pool := testPool(t)
	repo := NewUserRepo(pool)

	email := uniqueName("bob") + "@example.com"
	username := uniqueName("bob")
	createTestUser(t, repo, pool, email, username)

	tests := []struct {
		name     string
		email    string
		username string
		wantErr  error
	}{
		{"same email", email, uniqueName("other"), ErrEmailTaken},
		{"same email, different case", strings.ToUpper(email), uniqueName("other"), ErrEmailTaken},
		{"same username", uniqueName("other") + "@example.com", username, ErrUsernameTaken},
		{"same username, different case", uniqueName("other") + "@example.com", strings.ToUpper(username), ErrUsernameTaken},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.Create(context.Background(), tt.email, tt.username, "hash-not-real")
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Create error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepoGetByIDNotFound(t *testing.T) {
	pool := testPool(t)
	repo := NewUserRepo(pool)

	_, err := repo.GetByID(context.Background(), "00000000-0000-0000-0000-000000000000")
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("GetByID error = %v, want %v", err, ErrUserNotFound)
	}
}
