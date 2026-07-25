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
`GoalWeight`, `WeightUnit`, `GoalReps`, `Block`, `SortOrder`, `IsTimeBased`, `GoalSeconds`, `LinkedToNext`,
`CircuitID` and `WorkSeconds`. The default block is `"main"` when blank.

`SessionCircuit` is a named, ordered group of exercises within one session (`Rounds`, `TransitionSeconds`,
`SortOrder`). It is **copied** from a template circuit, never linked to one: a session records what was performed, so
editing the template afterwards must not rewrite it.

**Circuit membership is `CircuitID != nil` and nothing else.** A loose exercise is permitted to carry a non-zero
`WorkSeconds` — the repository stores what it is given — so reading a duration as membership would quietly pull an
exercise into a circuit it was never in.

`LinkedToNext` means "do not rest after me — go straight to the next exercise", which is how supersets are recorded. It
belongs to one exercise rather than to a pair, so reordering or deleting cannot orphan a group. Never render from it
directly: the controller derives an *effective* link (`SupersetLinked`) that also requires a next exercise in the same
block and a run of at most four.

Related types: `SessionSet` records each logged set (actual weight, reps, seconds, activity type, set number);
`CardioLog` records cardio intervals. `SessionExerciseView` is the read projection bundling an exercise with its sets
and cardio logs for template rendering, including `HitMax`, `BelowGoal`, `GoalRepMin`, `GoalRepMax`, `SupersetLinked`,
and `SupersetLabel` computed by the controller.

Writes take a `SessionExerciseInput` struct rather than positional arguments: `Create` already had nine, and circuit
membership would have made eleven, almost all `int` and `bool`, where a transposed pair still compiles. `CircuitID` is
a `*int64` rather than an index with a sentinel, so a literal that omits it means "loose", not "the first circuit".

`CreateBody(sessionID, circuits, exercises)` copies a whole template body in **one transaction**: circuits first, then
exercises, with each exercise's `CircuitID` — which arrives holding the *template* circuit id — rewritten to the id of
the session circuit copied from it. An exercise naming a circuit that was not submitted alongside it becomes loose,
because storing the dangling id would fail the foreign key and roll back the entire workout over one bad row.
`GetCircuitsBySession` reads a session's circuits in sort order.

The repository provides full CRUD plus `LogSet`, `CountSetsByExercise`, `DeleteSet` (with set-number re-sequencing),
`LogCardio`, `DeleteCardioLog`, `UpdateSortOrders`, `UpdateName`, and `UpdateLinkedToNext(id int64, linked bool)`. Like
the other write methods here it is *not* userID-scoped — ownership is enforced in the session controller. The session
controller is the primary consumer.
