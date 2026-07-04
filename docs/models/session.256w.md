---
level: 256w
parent: ../models.32w.md
relates-to:
  - program.128w.md
  - session_exercise.128w.md
  - ../controllers/session.128w.md
source: models/session.go, models/session_repository.go
---

# Session (full reference)

## Struct fields

| Field         | Go type   | ORM tag / notes                 |
|---------------|-----------|---------------------------------|
| ID            | int64     | `auto;pk`                       |
| ProgramID     | int64     | FK to programs.id               |
| UserID        | int64     | FK to users.id                  |
| PhaseNumber   | int       | 1-indexed phase within program  |
| WeekNumber    | int       | 1-indexed week within phase     |
| WorkoutNumber | int       | 1-indexed workout within week   |
| IsDeload      | bool      | true when week == weeksPerPhase |
| Date          | time.Time | `type(date)` — date only        |
| CreatedAt     | time.Time | `auto_now_add`                  |
| UpdatedAt     | time.Time | `auto_now`                      |

## Repository interface (SessionRepository)

```go
Create(programID, userID int64, phaseNumber, weekNumber, workoutNumber int, isDeload bool, date time.Time) (*Session, error)
CountByProgram(programID int64) (int, error)
LatestByProgram(programID int64) (*Session, error)
GetByID(id, userID int64) (*Session, error)
GetByProgram(programID int64) ([]*Session, error)
GetRecentByUser(userID int64, limit int) ([]*RecentSession, error)
Delete(id, userID int64) error
```

## Scheduling helpers (pure functions)

- `CalculatePhaseAndWeek(now, startDate time.Time, weeksPerPhase int)` — derives phase/week from calendar date; defaults
  `weeksPerPhase` to 8 if <= 0.
- `CalculateNextSession(sessionCount, weeksPerPhase, workoutsPerWeek int)` — computes next position from session count
  alone.
- `CalculateNextSessionFromLast(last *Session, weeksPerPhase, workoutsPerWeek int)` — increments workout, wrapping to
  next week/phase as needed.
- Deload flag: `isDeload = week == weeksPerPhase` in all three helpers.

## RecentSession projection

```go
type RecentSession struct {
    ID, ProgramID int64
    PhaseNumber, WeekNumber, WorkoutNumber int
    IsDeload bool
    Date time.Time
    ProgramName string  // joined from programs table
}
```

Retrieved via raw SQL join:
`SELECT s.*, p.name AS program_name FROM sessions s JOIN programs p ON p.id = s.program_id WHERE s.user_id = ? ORDER BY s.date DESC, s.id DESC LIMIT ?`.

## ORM / SQL patterns

- `GetByProgram`: `QueryTable.Filter("ProgramID").OrderBy("-Date", "-ID")`.
- `LatestByProgram`: same filter + `Limit(1).One()`; returns `nil, nil` on `ErrNoRows`.
- `Delete`: reads by ID first, checks `UserID`, then deletes — enforces ownership.
- `GetByID`: reads, then ownership check; returns `"not found"` on mismatch.

## Relationships

- Belongs to `Program` and `User`.
- Has many `SessionExercise` rows.
