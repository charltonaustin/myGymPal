---
level: 128w
parent: ../models.32w.md
deeper: session.256w.md
relates-to:
  - program.128w.md
  - session_exercise.128w.md
  - ../controllers/session.128w.md
source: models/session.go, models/session_repository.go
---

# Session

`Session` records a single workout instance within a `Program`. It tracks `ProgramID`, `UserID`, `PhaseNumber`,
`WeekNumber`, `WorkoutNumber`, `IsDeload`, and `Date`. Auto-managed `CreatedAt`/`UpdatedAt` timestamps are also stored.

Three pure functions assist scheduling: `CalculatePhaseAndWeek` derives the current phase/week from start date;
`CalculateNextSession` computes the next position from a session count; `CalculateNextSessionFromLast` increments from
the most recent session's values directly. The deload flag is automatically set when `week == weeksPerPhase`.

`RecentSession` is a join projection bundling session fields with `ProgramName` for dashboard display.

The `SessionRepository` provides `Create`, `CountByProgram`, `LatestByProgram`, `GetByID`, `GetByProgram`,
`GetRecentByUser`, and `Delete`. All reads are user-scoped. The session controller and dashboard controller are the
primary consumers.
