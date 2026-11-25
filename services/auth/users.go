package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// User is a row in the users table.
type User struct {
	ID           string
	Email        string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

// These are errors the rest of the service can check for. Callers compare
// with errors.Is(err, ErrEmailTaken) and never need to know about SQL.
var (
	ErrEmailTaken    = errors.New("email already registered")
	ErrUsernameTaken = errors.New("username already taken")
	ErrUserNotFound  = errors.New("user not found")
)

// UserRepo is the only code that runs SQL against the users table. Everything
// else calls its methods, like Create, so the SQL stays in one place.
type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

// Create inserts a new user and returns it with its generated ID. The
// password must already be hashed. The repo never sees the real password.
func (r *UserRepo) Create(ctx context.Context, email, username, passwordHash string) (User, error) {
	u := User{Email: email, Username: username, PasswordHash: passwordHash}

	// $1, $2, $3 are placeholders. pgx sends the values separately from the SQL,
	// which prevents SQL injection. Never build SQL by joining strings.
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (email, username, password_hash)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at`,
		email, username, passwordHash,
	).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		// Postgres code 23505 means "unique violation". The constraint name
		// says which of our two unique indexes was hit.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_email_lower_idx":
				return User{}, ErrEmailTaken
			case "users_username_lower_idx":
				return User{}, ErrUsernameTaken
			}
		}
		return User{}, fmt.Errorf("insert user: %w", err)
	}
	return u, nil
}

// GetByID looks a user up by their UUID.
func (r *UserRepo) GetByID(ctx context.Context, id string) (User, error) {
	var u User
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, username, password_hash, created_at
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}
