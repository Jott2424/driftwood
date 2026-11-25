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
