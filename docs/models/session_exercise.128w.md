---
level: 128w
parent: ../models.32w.md
deeper: session_exercise.256w.md
relates-to:
  - session.128w.md
  - exercise.128w.md
  - cardio_log.128w.md
  - ../controllers/session.128w.md
source: models/session_exercise.go, models/session_exercise_repository.go
---

# SessionExercise

`SessionExercise` captures an exercise entry within a session: `Name` (lowercase/trimmed on insert), `IsBodyweight`,
`GoalWeight`, `WeightUnit`, `GoalReps`, `Block`, `SortOrder`, `IsTimeBased`, and `GoalSeconds`. The default block is
`"main"` when blank.

Related types: `SessionSet` records each logged set (actual weight, reps, seconds, activity type, set number);
`CardioLog` records cardio intervals. `SessionExerciseView` is the read projection bundling an exercise with its sets
and cardio logs for template rendering, including `HitMax`, `BelowGoal`, `GoalRepMin`, and `GoalRepMax` flags computed
by the controller.

The repository provides full CRUD plus `LogSet`, `CountSetsByExercise`, `DeleteSet` (with set-number re-sequencing),
`LogCardio`, `DeleteCardioLog`, and `UpdateSortOrders`. The session controller is the primary consumer.
