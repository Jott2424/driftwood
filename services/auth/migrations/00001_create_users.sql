-- +goose Up
CREATE TABLE users (
    id            uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    email         text        NOT NULL,
    username      text        NOT NULL,
    password_hash text        NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now()
);

-- Unique on the lowercase value, so "Ann@x.com" and "ann@x.com" count as the
-- same email, and likewise for usernames.
CREATE UNIQUE INDEX users_email_lower_idx    ON users (lower(email));
CREATE UNIQUE INDEX users_username_lower_idx ON users (lower(username));

-- +goose Down
DROP TABLE users;
