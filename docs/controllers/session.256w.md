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
| POST   | /sessions/:id/exercises/:eid/delete             | DeleteExercise    |
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
- `c.Data["ExerciseBlocks"]` = `[]sessionExerciseBlock` (grouped by main/abs/cardio/stretch)
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
- `Sessions.Delete` — Delete

## AJAX endpoints

- **LogSet** — returns `{"id", "set_number", "is_time_based", "actual_seconds", "activity_type"}` when
  `X-Requested-With: XMLHttpRequest`
- **DeleteSet** — returns `{"status": "ok"}` on AJAX
- **ReorderExercises** — always returns JSON `{"ok": "1"}` or `{"error": "..."}`

## Key invariants

- After 3rd set: if `actualReps >= exercise.GoalReps` and `actualWeight >= convertedGoalWeight` (or goalWeight is 0),
  calls `Exercises.UpdateGoalWeight`.
- `isDeload` is true when `weekNumber == program.WeeksPerPhase`.
- HitMax: prev session's first N sets (N = phaseDefaultSets) all at or above repMax and goalWeight.
- Block ordering: main → abs → cardio → stretch.

## Flash messages

- `flash.Success("Workout deleted.")` — Delete success

## Relationship to other controllers

Reads programs from ProgramController's domain; uses exercise library maintained by ExerciseController; uses templates
from TemplateController.
