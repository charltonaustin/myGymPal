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
`GoalWeight`, `WeightUnit`, `GoalReps`, `Block`, `SortOrder`, `IsTimeBased`, `GoalSeconds`, and `LinkedToNext`. The
default block is `"main"` when blank.

`LinkedToNext` means "do not rest after me — go straight to the next exercise", which is how supersets are recorded. It
belongs to one exercise rather than to a pair, so reordering or deleting cannot orphan a group. Never render from it
directly: the controller derives an *effective* link (`SupersetLinked`) that also requires a next exercise in the same
block and a run of at most four.

Related types: `SessionSet` records each logged set (actual weight, reps, seconds, activity type, set number);
`CardioLog` records cardio intervals. `SessionExerciseView` is the read projection bundling an exercise with its sets
and cardio logs for template rendering, including `HitMax`, `BelowGoal`, `GoalRepMin`, `GoalRepMax`, `SupersetLinked`,
and `SupersetLabel` computed by the controller.

The repository provides full CRUD plus `LogSet`, `CountSetsByExercise`, `DeleteSet` (with set-number re-sequencing),
`LogCardio`, `DeleteCardioLog`, `UpdateSortOrders`, `UpdateName`, and `UpdateLinkedToNext(id int64, linked bool)`. Like
the other write methods here it is *not* userID-scoped — ownership is enforced in the session controller. The session
controller is the primary consumer.
