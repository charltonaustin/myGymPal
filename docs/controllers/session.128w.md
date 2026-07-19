---
level: 128w
parent: ../controllers.32w.md
deeper: session.256w.md
relates-to:
  - program.128w.md
  - exercise.128w.md
  - template.128w.md
source: controllers/session.go
---

The session controller is the most complex in the application, managing the full workout-logging lifecycle. GET
/programs/:id/sessions/new computes the next phase/week/workout numbers (supporting both sequential and latest-session
modes) and renders a session creation form with available templates. POST /programs/:id/sessions creates the session,
optionally copies template exercises into it (resolving goal weights and reps from the exercise library and phase
config), then redirects to the active session view.

GET /sessions/:id is the workout page; it enriches session exercises with current library goals (including each
exercise's preferred unit), converts set weights to each exercise's own unit, computes HitMax and BelowGoal flags by
comparing previous-session performance (with proper per-exercise unit conversion) against the current phase's rep max,
and embeds exercise library JSON for the add-exercise autocomplete.

GET /sessions/:id also computes supersets. `groupSessionExercises` walks each block's sort-ordered exercises and fills
in `SupersetLinked` (the effective link) and `SupersetLabel` (`A1`, `A2`, …) on each view, capping a run at four
members and ignoring a link on an exercise that is last in its block.

POST handlers manage: adding exercises, logging sets (with AJAX support returning JSON for the set row), logging and
deleting cardio activities, deleting exercises, deleting sets (AJAX), reordering exercises (AJAX), updating a single
exercise's display unit (`UpdateExerciseUnit`, AJAX), and toggling a superset chain (`UpdateLink`, AJAX). After the 3rd
set at or above goal reps and goal weight, the exercise library goal weight is automatically promoted via
`Exercises.UpdateGoalWeight`.

`UpdateLink` (POST /sessions/:id/exercises/:eid/link, form field `linked=true|false`) is a JSON endpoint: 401 when
unauthenticated, 404 when the session is not the caller's or the exercise is not the session's, and 400 when turning a
link on for an exercise that is last in its block or that would build a run of more than four. Turning a link off is
always allowed, which is what lets a stale link be cleared.
