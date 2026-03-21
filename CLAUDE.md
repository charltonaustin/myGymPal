# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

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

## Architecture

**Stack:** Go + Beego MVC, PostgreSQL via Beego ORM, golang-migrate, Bootstrap 5.3.3

**Repository pattern:**

- Interfaces in `models/interfaces.go`
- ORM implementations in `models/*_repository.go`
- Package-level repository variables declared in `controllers/repos.go`
- Repositories instantiated and injected in `routers/router.go`:`Register()` before `beego.Run()`
- Tests override these globals with mocks (see `controllers/mocks_test.go`)

**Request flow:** `routers/router.go` → `controllers/*.go` → `models/*_repository.go` → PostgreSQL

**Templates:** `.tpl` files in `views/`, data passed via `c.Data["key"]`, rendered with `c.TplName`. Custom template
functions (e.g., `fmtDuration`) registered in `routers/router.go`.

**AJAX responses:** `c.Data["json"] = payload` + `c.ServeJSON()`

**Auth:** Session-based (`c.GetSession("user_id")`), 30-day lifetime backed by PostgreSQL in production, in-memory for
tests.

**Flash messages:** `beego.NewFlash()` → `flash.Success(...)` → `flash.Store(c.Ctx)` →
`beego.ReadFromRequest(&c.Controller)`

## Key Invariants

- Exercise names are **always lowercase + trimmed** (enforced in model layer). `GetExerciseByName` uses
  `LOWER(TRIM(...))` on both sides.
- `template.JS` type is used for JSON embedded in `<script>` tags to prevent html/template escaping.
- Template exercises store only name + is_bodyweight; goal weight/reps come from the Exercise library and Phase config
  at render time.
- Auto-updates exercise goal weight after 3rd set at or above goal reps.

## Configuration

- `conf/app.conf` — development (PostgreSQL sessions, `mygympal` DB)
- `conf/app.test.conf` — tests (in-memory sessions, `mygympal_test` DB)

## Testing

Controller tests use HTTP test helpers (`getPath`, `postForm`) and mock repositories with function fields (`CreateFn`,
`GetByIDFn`, etc.) that can be overridden per-test. Model tests hit the real `mygympal_test` database.

`TestMain` in both `controllers/setup_test.go` and `models/user_test.go` handles Beego initialization and config
loading — no manual setup needed.
