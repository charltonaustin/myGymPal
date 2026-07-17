---
level: 32w
parent: README.16w.md
deeper: models/user.128w.md
relates-to:
  - controllers.32w.md
source: models/
---

| Model            | Summary                                                                                                                                                      | Detail                                  |
|------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------|
| user             | Authenticated user with username, bcrypt password hash, and weight unit preference (lb/kg); used by auth and account controllers.                            | [128w](models/user.128w.md)             |
| program          | A named training program owned by a user, with start date, phase count, weeks per phase, and workouts per week; creates phases transactionally on creation.  | [128w](models/program.128w.md)          |
| phase            | A numbered sub-period of a program storing rep range (min/max), default set count, and rest seconds; updated in bulk via PhaseUpdate. `PhaseRepository.UpdateRestSeconds(programID, phaseNumber, restSeconds)` updates only the rest period for a single phase (backs the in-session rest-timer control). | [128w](models/phase.128w.md)            |
| session          | A single workout instance within a program, tagged with phase, week, and workout numbers plus a deload flag; supports next-session calculation helpers. `SessionRepository.GetDailyActivity(userID, days)` returns `[]DayActivity` (one per exercised day within the window, with that day's logged-set count) for the activity heatmap.      | [128w](models/session.128w.md)          |
| session_exercise | An exercise entry within a session with goal weight/reps/seconds, block assignment, and sort order; paired with SessionSet and CardioLog sub-records.        | [128w](models/session_exercise.128w.md) |
| template         | A reusable workout template with name, focus, and an ordered list of named exercises (body-weight or time-based); created and updated transactionally.       | [128w](models/template.128w.md)         |
| program_workout_template | Per-program default template for each workout number (1–workouts_per_week); stored in `program_workout_templates` with upsert semantics; used to pre-select the template when starting a workout. | — |
| exercise         | Global exercise name (shared across users) in `exercises` table; per-user goals/config (weight, reps, seconds, block, bodyweight/time-based flags) in `user_exercise_goals`. `ExerciseRepository.GetAll` returns all global exercises with user goals overlaid; `GetAllByUser` returns only user-configured ones; `GetAvailableGlobalNames` returns names of global exercises the user has no goal for yet (used for new-exercise autocomplete). `GetHistory(userID, names, unit, days)` returns `[]ExerciseHistorySeries` — per-session max weight (normalized to `unit`) for standard exercises and max reps for bodyweight exercises; time-based exercises are skipped; only sessions within the last `days` days are included (`days <= 0` means all history). `GetRecentNames(userID, days)` returns names of exercises the user logged a set for within the last `days` days (most-recently-performed first) to pre-populate the history graph. | [128w](models/exercise.128w.md) |
| macro_entry      | A single food log entry per user per date recording food name, serving size, protein, carbs, and fat macros; supports distinct-food lookup for autocomplete. | [128w](models/macro_entry.128w.md)      |
| macro_goal       | A single per-user daily macro target (protein, carbs, fat grams) stored with upsert semantics; one row per user with auto-updated timestamp.                 | [128w](models/macro_goal.128w.md)       |
| body_weight      | A dated body-weight measurement per user storing weight value and unit; ordered newest-first for charting and tracking purposes.                             | [128w](models/body_weight.128w.md)      |
| cardio_log       | A cardio activity sub-record linked to a session exercise, storing cardio type, goal duration, and actual duration in seconds.                               | [128w](models/cardio_log.128w.md)       |
| interfaces       | Canonical Go interfaces for all repository types; the source of truth for mock generation and dependency injection in controllers and tests.                 | [128w](models/interfaces.128w.md)       |
