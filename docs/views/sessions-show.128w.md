---
level: 128w
parent: sessions.32w.md
deeper: sessions-show.256w.md
relates-to:
  - ../controllers/sessions.128w.md
source: views/sessions/show.tpl
---

Live workout session page. Includes navbar partial, exercise-fields partial (for the "Add Exercise" form), and all three
goal modals. Renders exercise blocks (main, abs, cardio, stretch) with sets tables. Set-logging forms submit via AJAX (
`fetch`) and append rows without reload; delete-set forms also use AJAX. Exercises in non-cardio blocks are
drag-to-reorder via SortableJS; reorder fires `POST /sessions/:id/exercises/reorder`. A fixed-bottom rest timer starts on each set logged, uses Web Audio API for alarm, and persists across page reloads via
`localStorage`. A global lb/kg toggle converts all exercises at once; each non-time-based exercise card also has its own
per-exercise lb/kg toggle that persists the preference to the exercise library via
`POST /sessions/:id/exercises/:eid/unit`. Each exercise card header has a three-dots (⋮) dropdown button — icon shows ↑/↓ for HitMax/BelowGoal — whose menu
contains "Edit goal" (opens the appropriate goal modal) and "Change exercise" (modal with autocomplete that POSTs
`name` to `/sessions/:id/exercises/:eid/change` and reloads the page). Exercise autocomplete fills type radio from `.ExerciseLibraryJSON`.
