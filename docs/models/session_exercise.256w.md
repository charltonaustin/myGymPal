---
level: 256w
parent: ../models.32w.md
relates-to:
  - session.128w.md
  - exercise.128w.md
  - cardio_log.128w.md
  - ../controllers/session.128w.md
source: models/session_exercise.go, models/session_exercise_repository.go
---

# SessionExercise (full reference)

## Struct fields — SessionExercise

| Field        | Go type | Notes                                                            |
|--------------|---------|------------------------------------------------------------------|
| ID           | int64   | `auto;pk`                                                        |
| SessionID    | int64   | FK to sessions.id                                                |
| Name         | string  | always lowercase+trimmed on insert                               |
| IsBodyweight | bool    |                                                                  |
| GoalWeight   | float64 |                                                                  |
| WeightUnit   | string  | `"lb"` or `"kg"`                                                 |
| GoalReps     | int     |                                                                  |
| Block        | string  | `"main"`, `"abs"`, `"cardio"`, `"stretch"`; defaults to `"main"` |
| SortOrder    | int     | set to current count at insert time                              |
| IsTimeBased  | bool    |                                                                  |
| GoalSeconds  | int     |                                                                  |

## Struct fields — SessionSet

| Field             | Go type | Notes                             |
|-------------------|---------|-----------------------------------|
| ID                | int64   | `auto;pk`                         |
| SessionExerciseID | int64   | FK to session_exercises.id        |
| SetNumber         | int     | 1-indexed; re-sequenced on delete |
| ActualWeight      | float64 |                                   |
| WeightUnit        | string  |                                   |
| ActualReps        | int     |                                   |
| ActualSeconds     | int     |                                   |
| ActivityType      | string  |                                   |

`SessionSet` helper methods: `Hours()`, `Minutes()`, `Secs()` decompose `ActualSeconds`.

## SessionExerciseView (read projection)

```go
type SessionExerciseView struct {
    Exercise   *SessionExercise
    Sets       []*SessionSet
    CardioLogs []*CardioLog
    HitMax     bool
    BelowGoal  bool
    GoalRepMin int
    GoalRepMax int
}
```

`LastSet()` returns the last element of `Sets` or `nil`.

## Repository interface (SessionExerciseRepository)

```go
Create(sessionID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, goalReps int, block string, isTimeBased bool, goalSeconds int) (*SessionExercise, error)
GetBySession(sessionID int64) ([]*SessionExerciseView, error)
GetByID(exerciseID int64) (*SessionExercise, error)
LogSet(exerciseID int64, setNumber int, actualWeight float64, weightUnit string, actualReps int, actualSeconds int, activityType string) (*SessionSet, error)
CountSetsByExercise(exerciseID int64) (int, error)
DeleteSet(setID int64) error
LogCardio(sessionExerciseID int64, cardioType string, goalDuration, actualDuration int) (*CardioLog, error)
DeleteCardioLog(id int64) error
DeleteExercise(exerciseID int64) error
UpdateSortOrders(sessionID int64, ids []int64) error
```

## Notable behavior

- `Create` normalizes name with `strings.ToLower(strings.TrimSpace(...))` and defaults block to `"main"`.
- `SortOrder` is set to the current row count for the session at insert time.
- `DeleteSet` uses a transaction: deletes the row, then decrements `set_number` for all later sets of the same
  exercise (raw SQL).
- `GetBySession` loads sets ordered by `SetNumber` and cardio logs ordered by `CreatedAt` for `"cardio"` block exercises
  only.
- `UpdateSortOrders` issues one raw `UPDATE` per id in the provided slice.
