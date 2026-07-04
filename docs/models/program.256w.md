---
level: 256w
parent: ../models.32w.md
relates-to:
  - ../controllers/program.128w.md
  - phase.128w.md
  - session.128w.md
source: models/program.go, models/program_repository.go
---

# Program (full reference)

## Struct fields

| Field           | Go type   | ORM tag / notes                             |
|-----------------|-----------|---------------------------------------------|
| ID              | int64     | `auto;pk`                                   |
| UserID          | int64     | FK to users.id                              |
| Name            | string    | non-empty enforced                          |
| StartDate       | time.Time | `type(date)` — date only, no time component |
| NumPhases       | int       | must be > 0                                 |
| WeeksPerPhase   | int       | must be > 0                                 |
| WorkoutsPerWeek | int       | must be > 0                                 |
| CreatedAt       | time.Time | `auto_now_add`                              |
| UpdatedAt       | time.Time | `auto_now`                                  |

## Repository interface (ProgramRepository)

```go
Create(userID int64, name string, startDate time.Time, numPhases, weeksPerPhase, workoutsPerWeek, defaultRepMin, defaultRepMax, defaultSets int) (*Program, error)
GetAllByUser(userID int64) ([]*Program, error)
GetByID(id, userID int64) (*Program, error)
Delete(id, userID int64) error
```

## Notable behavior

- `Create` runs inside a transaction: inserts the `Program`, then inserts `numPhases` `Phase` rows using
  `defaultRepMin`, `defaultRepMax`, and `defaultSets`. Rolls back on any failure.
- `GetAllByUser` orders results by `StartDate` ascending.
- `GetByID` reads by primary key, then checks `p.UserID == userID`; returns `"not found"` if ownership fails (no 401/403
  distinction in model layer).
- `Delete` filters by both `ID` and `UserID` and returns `"not found"` if zero rows are affected.

## Validation rules

- `name` must be non-empty.
- `numPhases`, `weeksPerPhase`, `workoutsPerWeek`, `defaultRepMin` must all be > 0.
- `defaultRepMax >= defaultRepMin`.

## ORM / SQL patterns

- `orm.RegisterModel(&Program{})` in `init()`.
- `GetAllByUser`: `QueryTable.Filter("UserID", userID).OrderBy("StartDate")`.
- `Delete`: `QueryTable.Filter("ID").Filter("UserID").Delete()`.
- Transactional creation uses `orm.NewOrm().Begin()` / `tx.Commit()` / `tx.Rollback()`.

## Relationships

- Belongs to `User` via `UserID`.
- Has many `Phase` rows (one per phase, created on program creation).
- Has many `Session` rows (created later as workouts are logged).
