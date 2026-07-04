---
level: 128w
parent: ../models.32w.md
deeper: exercise.256w.md
relates-to:
  - session_exercise.128w.md
  - template.128w.md
  - ../controllers/exercise.128w.md
source: models/exercise.go, models/exercise_repository.go
---

# Exercise

`Exercise` is the per-user exercise library. Each entry stores a lowercase-trimmed `Name`, `IsBodyweight`, `GoalWeight`,
`WeightUnit`, `IsTimeBased`, `GoalSeconds`, `GoalRepMin`, `GoalRepMax`, and `DefaultBlock`. Auto-managed timestamps are
included.

Name normalization (`strings.ToLower(strings.TrimSpace(...))`) is enforced on both `Create` and `Update`. `GetByName`
uses a raw SQL query with `LOWER(TRIM(...))` on both sides for case-insensitive lookup. `DefaultBlock` is validated
against `"abs"`, `"cardio"`, and `"stretch"`; any other value falls back to `"main"`.

`UpdateGoalWeight` provides a targeted update of only `GoalWeight`, `WeightUnit`, and `UpdatedAt` — used for the
auto-promotion after 3 successful sets.

The `ExerciseRepository` interface provides full CRUD plus `GetByName` and `UpdateGoalWeight`. The exercise controller
and session controller are the primary consumers.
