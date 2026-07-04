# CLAUDE.md

## Commands

```bash
# Start the database (required before running the app or tests)
docker compose up -d

# Run the application
go run main.go   # listens on localhost:8080

# Run all tests (database must be running)
go test ./...
go test ./... -v

# Run tests in a specific package
go test ./models/... -v
go test ./controllers/... -v

# Run a single test
go test ./models -run TestCreateUser_Success -v
go test ./controllers -run TestRegisterPost_Success -v
```

Migrations run automatically on startup. No manual migration commands needed.

## Docs

Architecture: see [docs/README.16w.md](docs/README.16w.md)
Key invariants: see [docs/models.32w.md](docs/models.32w.md)
Configuration: see [docs/conf.32w.md](docs/conf.32w.md)
Testing patterns: see [docs/controllers.32w.md](docs/controllers.32w.md)
