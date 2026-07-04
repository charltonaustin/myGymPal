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

GET /sessions/:id is the workout page; it enriches session exercises with current library goals, converts weights to the
user's preferred unit, computes HitMax and BelowGoal flags by comparing previous-session performance against the current
phase's rep max, and embeds exercise library JSON for the add-exercise autocomplete.

POST handlers manage: adding exercises, logging sets (with AJAX support returning JSON for the set row), logging and
deleting cardio activities, deleting exercises, deleting sets (AJAX), and reordering exercises (AJAX). After the 3rd set
at or above goal reps and goal weight, the exercise library goal weight is automatically promoted via
`Exercises.UpdateGoalWeight`.
