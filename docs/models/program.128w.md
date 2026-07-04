---
level: 128w
parent: ../models.32w.md
deeper: program.256w.md
relates-to:
  - ../controllers/program.128w.md
  - phase.128w.md
  - session.128w.md
source: models/program.go, models/program_repository.go
---

# Program

`Program` is a named training plan owned by a user. It records the program's `Name`, `StartDate`, `NumPhases`,
`WeeksPerPhase`, and `WorkoutsPerWeek`, plus auto-managed `CreatedAt`/`UpdatedAt` timestamps.

Creation is transactional: `CreateProgram` inserts the program row and then inserts one `Phase` row per phase (using
`defaultRepMin`, `defaultRepMax`, and `defaultSets`) in a single database transaction. This guarantees phases always
exist when a program is created.

Validation on creation enforces that `name` is non-empty, `numPhases`, `weeksPerPhase`, and `workoutsPerWeek` are all
positive, and `defaultRepMax >= defaultRepMin`.

The `ProgramRepository` interface provides `Create`, `GetAllByUser` (ordered by `StartDate`), `GetByID` (with user
ownership check), and `Delete`. The program controller uses all four methods; sessions and phases are retrieved by
program ID.
