# My Gym Pal

A mobile-friendly web app for tracking gym workouts. Log exercises, track weight progression, and manage training
programs across phases and weeks.

Built with Go, Beego v2, and PostgreSQL.

Try it at https://my-gym-pal.org/

---

## Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [Docker](https://www.docker.com/) (for the database)

---

## Setup

### 1. Install dependencies

```bash
go mod download
```

### 2. Start the database

```bash
docker compose up -d
```

This starts a PostgreSQL 17 container on port `5432` with the `mygympal` database. Credentials are in `conf/app.conf`.

### 3. Run the application

```bash
go run main.go
```

Database migrations run automatically on startup. The server listens at `http://localhost:8080`.

---

## Running Tests

The test suite requires the database to be running.

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run a specific package
go test ./models/... -v
go test ./controllers/... -v
```

Tests run against a separate `mygympal_test` database (configured in `conf/app.test.conf`) so they never touch
development data. The test database is created automatically when the Docker container first starts.

---

## Configuration

All configuration lives in `conf/app.conf`. Key settings:

| Key           | Default     | Description         |
|---------------|-------------|---------------------|
| `httpaddr`    | `localhost` | Server bind address |
| `httpport`    | `8080`      | Server port         |
| `runmode`     | `dev`       | `dev` or `prod`     |
| `db_host`     | `localhost` | PostgreSQL host     |
| `db_port`     | `5432`      | PostgreSQL port     |
| `db_name`     | `mygympal`  | Database name       |
| `db_user`     | `postgres`  | Database user       |
| `db_password` | `postgres`  | Database password   |

---

## Database

Connect to the running PostgreSQL container:

```bash
docker exec -it mygympal_db psql -U postgres -d mygympal
```

Stop the container (data is persisted in a Docker volume):

```bash
docker compose down
```

Wipe all data:

```bash
docker compose down -v
```
