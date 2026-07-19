---
level: 256w
parent: ../controllers.32w.md
relates-to:
  - program.256w.md
  - exercise.256w.md
  - template.256w.md
source: controllers/session.go
---

## Routes

| Method | Path                                            | Handler           |
|--------|-------------------------------------------------|-------------------|
| GET    | /programs/:id/sessions/new                      | New               |
| POST   | /programs/:id/sessions                          | Create            |
| GET    | /sessions/:id                                   | Show              |
| POST   | /sessions/:id/delete                            | Delete            |
| POST   | /sessions/:id/exercises                         | AddExercise       |
| POST   | /sessions/:id/exercises/reorder                 | ReorderExercises  |
| POST   | /sessions/:id/exercises/:eid/unit               | UpdateExerciseUnit |
| POST   | /sessions/:id/exercises/:eid/change             | ChangeExercise     |
| POST   | /sessions/:id/exercises/:eid/delete             | DeleteExercise    |
| POST   | /sessions/:id/exercises/:eid/link               | UpdateLink        |
| POST   | /sessions/:id/exercises/:eid/sets               | LogSet            |
| POST   | /sessions/:id/exercises/:eid/sets/:sid/delete   | DeleteSet         |
| POST   | /sessions/:id/cardio                            | AddCardioActivity |
| POST   | /sessions/:id/exercises/:eid/cardio             | LogCardio         |
| POST   | /sessions/:id/exercises/:eid/cardio/:lid/delete | DeleteCardioLog   |

## Auth requirement

All handlers check `c.GetSession("user_id")`; nil → /login (or JSON error for AJAX handlers).

## Session keys

- Read: `user_id` (int64)

## Template variables — New

- `c.Data["Program"]`, `c.Data["PhaseNumber"]`, `c.Data["WeekNumber"]`, `c.Data["WorkoutNumber"]`
- `c.Data["Templates"]`, `c.Data["DefaultDate"]`, `c.Data["LogMode"]`

## Template variables — Show

- `c.Data["Session"]`, `c.Data["Program"]`
- `c.Data["ExerciseBlocks"]` = `[]sessionExerciseBlock` (grouped by main/abs/cardio/stretch). Each
  `SessionExerciseView` inside carries the computed `SupersetLinked` and `SupersetLabel`; they are not separate
  top-level `c.Data` keys.
- `c.Data["WeightUnit"]`, `c.Data["ExWeightUnit"]`
- `c.Data["PhaseRepMin"]`, `c.Data["PhaseRepMax"]`, `c.Data["PhaseRestSeconds"]`
- `c.Data["ExerciseLibraryJSON"]` = `template.JS` — embedded JSON for autocomplete

## Templates

- `sessions/new.tpl`, `sessions/show.tpl`

## Repository calls

- `Programs.GetByID`, `Sessions.CountByProgram`, `Sessions.LatestByProgram` — New
- `Sessions.Create`, `Templates.GetByID`, `Phases.GetByProgram`, `Exercises.GetByName`, `SessionExercises.Create` —
  Create
- `Sessions.GetByID`, `Programs.GetByID`, `SessionExercises.GetBySession`, `Exercises.GetByName`, `Users.GetByID`,
  `Phases.GetByProgram`, `Sessions.GetByProgram`, `SessionExercises.GetBySession` (previous) — Show
- `SessionExercises.Create` — AddExercise, AddCardioActivity
- `SessionExercises.LogSet`, `SessionExercises.CountSetsByExercise`, `Exercises.GetByName`,
  `Exercises.UpdateGoalWeight` — LogSet
- `SessionExercises.LogCardio` — LogCardio, AddCardioActivity
- `SessionExercises.DeleteCardioLog` — DeleteCardioLog
- `SessionExercises.DeleteExercise` — DeleteExercise
- `SessionExercises.DeleteSet` — DeleteSet
- `SessionExercises.UpdateSortOrders` — ReorderExercises
- `Sessions.GetByID`, `SessionExercises.GetByID`, `Exercises.GetByName`, `Exercises.UpdateGoalWeight` — UpdateExerciseUnit
- `Sessions.GetByID`, `SessionExercises.GetByID`, `SessionExercises.UpdateName` — ChangeExercise
- `Sessions.GetByID`, `SessionExercises.GetByID`, `SessionExercises.GetBySession`,
  `SessionExercises.UpdateLinkedToNext` — UpdateLink
- `Sessions.Delete` — Delete

## AJAX endpoints

- **LogSet** — returns `{"id", "set_number", "is_time_based", "actual_seconds", "activity_type"}` when
  `X-Requested-With: XMLHttpRequest`
- **DeleteSet** — returns `{"status": "ok"}` on AJAX
- **ReorderExercises** — always returns JSON `{"ok": "1"}` or `{"error": "..."}`
- **UpdateExerciseUnit** — AJAX `POST /sessions/:id/exercises/:eid/unit` with `weight_unit=lb|kg`; looks up exercise
  library entry by session-exercise name, converts goal weight, saves new unit; returns `{"ok": true}`
- **ChangeExercise** — `POST /sessions/:id/exercises/:eid/change` with `name`; validates ownership, updates
  `session_exercises.name`, redirects back to session page; existing sets remain linked to the exercise row
- **UpdateLink** — `POST /sessions/:id/exercises/:eid/link` with `linked=true|false`; toggles the superset chain to the
  next exercise in the block. Returns `{"ok": true, "linked": <bool>}`; 401 unauthenticated, 404 on a session or
  exercise the caller does not own, 400 when linking an exercise that is last in its block or that would make a run of
  more than four. Never redirects — an AJAX caller cannot follow one.

## Key invariants

- After 3rd set: if `actualReps >= exercise.GoalReps` and `actualWeight >= convertedGoalWeight` (or goalWeight is 0),
  calls `Exercises.UpdateGoalWeight`.
- `isDeload` is true when `weekNumber == program.WeeksPerPhase`.
- HitMax: prev session's first N sets (N = phaseDefaultSets) all at or above repMax and goalWeight; prev-set weights are converted to the exercise's own unit before comparison.
- Block ordering: main → abs → cardio → stretch.
- Supersets: `linked_to_next` is *never* read directly at render time. `groupSessionExercises` computes
  `SupersetLinked` = raw link AND an exercise exists at `i+1` in the same block AND the run stays under four members;
  runs of two or more are labelled `A1`, `A2`, `B1`, … per block. A stale link on the last exercise of a block is
  inert, so the rest timer still fires after it. `views/sessions/show.tpl` starts the rest timer after a logged set
  only when the card's computed `data-linked` is not `true`.
- Ownership for `UpdateLink` is enforced in the controller (`Sessions.GetByID(sessionID, userID)`, then
  `exercise.SessionID != sessionID`); `SessionExerciseRepository` methods are not userID-scoped.

## Flash messages

- `flash.Success("Workout deleted.")` — Delete success

## Relationship to other controllers

Reads programs from ProgramController's domain; uses exercise library maintained by ExerciseController; uses templates
from TemplateController.
