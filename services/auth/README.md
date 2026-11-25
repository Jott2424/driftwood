# auth

Accounts, login, profiles, and token issuing.

## Run locally

Start the database (from the repo root):

```
docker compose -f deploy/compose.yaml up -d
```

Then run the service. PowerShell:

```
$env:DATABASE_URL = "postgres://driftwood:driftwood_dev@localhost:5432/auth"
go run .
```

| Endpoint | Meaning |
|---|---|
| `GET /health` | The process is alive. |
| `GET /ready` | The service can reach its database (503 if not). |

## Database migrations

The SQL files in `migrations/` are compiled into the program and applied
automatically at startup. goose records which ones have run in a table called
`goose_db_version`, so each runs once.

To change the database structure, add a new numbered file, like
`migrations/00002_add_something.sql`, with a `-- +goose Up` section and a
`-- +goose Down` section that undoes it. Never edit a migration that has already
been applied somewhere. Add a new one instead.

## Tests

```
go test ./...
```

Integration tests that need a real database are skipped unless
`TEST_DATABASE_URL` is set. To include them:

```
$env:TEST_DATABASE_URL = "postgres://driftwood:driftwood_dev@localhost:5432/auth"
go test ./...
```
