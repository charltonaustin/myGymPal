---
level: 256w
parent: ../controllers.32w.md
relates-to:
  - default.256w.md
  - auth.256w.md
source: controllers/error.go
---

## Routes

| Method | Path   | Handler |
|--------|--------|---------|
| GET    | /error | Get     |

## Auth requirement

None. No session check of any kind.

## Session keys

None read or written.

## Template variables

None. `c.Data` is not populated.

## Template

- `error.tpl`

## Repository calls

None.

## Flash messages

None.

## Redirect paths

This controller is itself a redirect target; it does not redirect to anything else.

## Usage across the codebase

`c.Redirect("/error", 302)` is called by:

- **AuthController.Logout** — if `c.DestroySession()` fails
- **AuthController.LoginPost** — if `c.SetSession` fails
- **ProgramController** — if `Programs.GetAllByUser` or `Phases.GetByProgram` fail
- **SessionController** — if `Sessions.CountByProgram` or `Sessions.LatestByProgram` fail; if `Sessions.Create` fails
- **ExerciseController** — if `Exercises.GetAllByUser` fails

## Relationship to other controllers

All authenticated controllers that encounter unrecoverable data-layer failures redirect here. The error page is
intentionally stateless so it is reachable even when sessions or DB connectivity are partially degraded.
