---
level: 256w
parent: ../models.32w.md
relates-to:
  - session_exercise.128w.md
  - ../controllers/session.128w.md
source: models/cardio_log.go
---

# CardioLog (full reference)

## Struct fields

| Field             | Go type   | ORM tag / notes                          |
|-------------------|-----------|------------------------------------------|
| ID                | int64     | `auto;pk`                                |
| SessionExerciseID | int64     | FK to session_exercises.id               |
| CardioType        | string    | free-form label (e.g. `"run"`, `"bike"`) |
| GoalDuration      | int       | target duration in seconds               |
| ActualDuration    | int       | actual duration in seconds               |
| CreatedAt         | time.Time | `auto_now_add`                           |

No `UpdatedAt` — cardio logs are append-only (no update path).

## Package-level functions

```go
LogCardioEntry(sessionExerciseID int64, cardioType string, goalDuration, actualDuration int) (*CardioLog, error)
GetCardioLogsByExercise(sessionExerciseID int64) ([]*CardioLog, error)
DeleteCardioLog(id int64) error
```

## Repository exposure

`CardioLog` has no dedicated repository interface. Operations are surfaced through `SessionExerciseRepository`:

- `LogCardio(sessionExerciseID int64, cardioType string, goalDuration, actualDuration int) (*CardioLog, error)` —
  delegates to `LogCardioEntry`.
- `DeleteCardioLog(id int64) error` — delegates to `DeleteCardioLog`.

`GetCardioLogsByExercise` is not in any interface; it is called internally by `GetSessionExercisesWithSets` when
building `SessionExerciseView` for `"cardio"` block exercises.

## Notable behavior

- `LogCardioEntry`: plain insert, no validation beyond ORM constraints.
- `GetCardioLogsByExercise`: `QueryTable.Filter("SessionExerciseID").OrderBy("CreatedAt")` — chronological order.
- `DeleteCardioLog`: deletes by primary key with `o.Delete(&CardioLog{ID: id})` — no ownership check (ownership is
  enforced at the session/session-exercise level by the controller).
- No update path exists — logs are immutable after creation.

## ORM / SQL patterns

- `orm.RegisterModel(&CardioLog{})` in `init()`.
- All operations use the default ORM instance; no raw SQL.

## Relationships

- Belongs to `SessionExercise` via `SessionExerciseID`.
- Only populated for exercises where `Block == "cardio"`.
- Embedded in `SessionExerciseView.CardioLogs []*CardioLog` for template rendering.
