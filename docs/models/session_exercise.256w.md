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
| LinkedToNext | bool    | raw superset flag: "no rest after me"; never render it directly  |
| CircuitID    | *int64  | nil for a loose exercise; **the only** marker of circuit membership |
| WorkSeconds  | int     | length of one round; meaningful only inside a circuit, but a loose exercise may legally carry one |

## Struct fields — SessionCircuit

| Field             | Go type | Notes                                     |
|-------------------|---------|-------------------------------------------|
| ID                | int64   | `auto;pk`                                 |
| SessionID         | int64   | FK to sessions.id, `ON DELETE CASCADE`    |
| Name              | string  |                                           |
| Rounds            | int     | `>= 1`, clamped by `ValidRounds`          |
| TransitionSeconds | int     | `>= 0`, clamped by `ValidSeconds`         |
| SortOrder         | int     |                                           |

A session circuit is a **copy** of a template circuit, never a reference to one. Editing a template next month must not
change the rounds or durations of a workout already performed.

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
    Exercise       *SessionExercise
    Sets           []*SessionSet
    CardioLogs     []*CardioLog
    HitMax         bool
    BelowGoal      bool
    GoalRepMin     int
    GoalRepMax     int
    SupersetLinked bool   // effective link, computed per block by the controller
    SupersetLabel  string // "A1", "A2", … or "" for a solo exercise
}
```

`LastSet()` returns the last element of `Sets` or `nil`.

## Superset link

`UpdateSessionExerciseLink(id int64, linked bool) error` writes `linked_to_next` by raw SQL (same idiom as
`UpdateSessionExerciseName`); the repository exposes it as `UpdateLinkedToNext`. The column arrives back through
`GetSessionExercisesWithSets` automatically, since the ORM selects the whole struct.

`SupersetLinked` and `SupersetLabel` are *computed*, never stored. `controllers.groupSessionExercises` walks each
block's sort-ordered slice and marks exercise `i` linked only when `LinkedToNext` is set, an exercise exists at `i+1`
in the same block, and the run it extends holds fewer than four members. Maximal runs of two or more get a per-block
letter and a 1-based index (`A1`, `A2`, `B1`, …); solo exercises get `""`. A stale `true` on the last exercise of a
block is therefore inert, and the rest timer fires after it.

## Repository interface (SessionExerciseRepository)

```go
Create(in SessionExerciseInput) (*SessionExercise, error)
CreateBody(sessionID int64, circuits []SessionCircuitInput, exercises []SessionExerciseInput) error
GetCircuitsBySession(sessionID int64) ([]*SessionCircuit, error)
GetBySession(sessionID int64) ([]*SessionExerciseView, error)
GetByID(exerciseID int64) (*SessionExercise, error)
LogSet(exerciseID int64, setNumber int, actualWeight float64, weightUnit string, actualReps int, actualSeconds int, activityType string) (*SessionSet, error)
CountSetsByExercise(exerciseID int64) (int, error)
DeleteSet(setID int64) error
LogCardio(sessionExerciseID int64, cardioType string, goalDuration, actualDuration int) (*CardioLog, error)
DeleteCardioLog(id int64) error
DeleteExercise(exerciseID int64) error
UpdateSortOrders(sessionID int64, ids []int64) error
UpdateName(id int64, name string) error
UpdateLinkedToNext(id int64, linked bool) error
```

## Input structs

```go
type SessionExerciseInput struct {
    SessionID    int64
    Name         string
    IsBodyweight bool
    GoalWeight   float64
    WeightUnit   string
    GoalReps     int
    Block        string
    IsTimeBased  bool
    GoalSeconds  int
    CircuitID    *int64  // nil = loose; pointer, so an omitted field is not "circuit 0"
    WorkSeconds  int
}

type SessionCircuitInput struct {
    TemplateCircuitID int64  // matching key only; never stored
    Name              string
    Rounds            int
    TransitionSeconds int
    SortOrder         int
}
```

`Create` took nine positional parameters, seven of them `int`, `bool` or `string`. Circuits needed two more, and an
eleven-argument call of near-identical types is where a transposed pair compiles cleanly, passes every test, and
surfaces months later as a workout with the wrong goal.

## Circuits

`CreateSessionBody(sessionID, circuits, exercises)` runs in one transaction:

1. insert every `SessionCircuit`, recording `TemplateCircuitID → new id`;
2. insert every exercise at `existingCount + i`, rewriting `CircuitID` from the template id to the session id.

An exercise whose `CircuitID` matches no submitted circuit is stored **loose**. The alternative — writing the dangling
template id — violates `session_exercises_circuit_id_fkey` and rolls back the whole copy, so one bad row would cost the
engineer the entire workout. (Confirmed by mutation: removing the fallback produces
`pq: ... violates foreign key constraint "session_exercises_circuit_id_fkey" (23503)`.)

`GetSessionCircuits(sessionID)` reads circuits ordered by `SortOrder`. Like the exercise reads it is scoped by session
alone; the *user* ownership check lives in the session controller (`Sessions.GetByID(id, userID)`).

A round is a set: the runner logs round N of an exercise as set N through the ordinary `LogSet` path, so no separate
results table exists and circuit work appears in the normal exercise history.

## Notable behavior

- `Create` normalizes name with `strings.ToLower(strings.TrimSpace(...))` and defaults block to `"main"`. It shares
  one row builder (`newSessionExercise`) with `CreateSessionBody`, so a field one path writes cannot go missing from
  the other — the failure mode that silently cleared `is_time_based` on every template edit.
- `SortOrder` is set to the current row count for the session at insert time.
- `DeleteSet` uses a transaction: deletes the row, then decrements `set_number` for all later sets of the same
  exercise (raw SQL).
- `GetBySession` loads sets ordered by `SetNumber` and cardio logs ordered by `CreatedAt` for `"cardio"` block exercises
  only.
- `UpdateSortOrders` issues one raw `UPDATE` per id in the provided slice.
