---
level: 128w
parent: ../models.32w.md
deeper: interfaces.256w.md
relates-to:
  - user.128w.md
  - program.128w.md
  - phase.128w.md
  - session.128w.md
  - session_exercise.128w.md
  - template.128w.md
  - exercise.128w.md
  - macro_goal.128w.md
  - body_weight.128w.md
source: models/interfaces.go
---

# Interfaces

`interfaces.go` is the canonical declaration file for all repository interfaces in the models package. It defines
`UserRepository`, `ProgramRepository`, `PhaseRepository`, `TemplateRepository`, `SessionRepository`,
`ExerciseRepository`, `BodyWeightRepository`, `MacroGoalRepository`, and `SessionExerciseRepository`.

These interfaces are the contracts used throughout the application: ORM implementations live in `*_repository.go` files,
mock implementations for tests live in `controllers/mocks_test.go`, and package-level variables in
`controllers/repos.go` hold the active instance.

Note: `MacroRepository` (for `MacroEntry`) is defined in `macro_entry_repository.go`, not here — it is the one
repository interface outside this file.

This file is the primary reference for understanding what each repository can do without reading implementation details.
