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
drag-to-reorder via SortableJS; reorder fires `POST /sessions/:id/exercises/reorder`. A fixed-bottom rest timer starts
on each set logged *unless the card is part of a superset*, uses Web Audio API for alarm, and persists across page
reloads via `localStorage`.

Supersets: every exercise card except the last in its block carries a chain button (`bi-link-45deg` off,
`bi-link` on) that toggles `POST /sessions/:id/exercises/:eid/link` and updates the card in place — the page never
reloads, so a running rest timer survives. Cards carry `data-link-raw` (the stored `linked_to_next`) and `data-linked`
(the computed effective link); members of a run show an `A1`/`A2` badge and share a left accent stripe. After a set is
logged, the timer starts only when `data-linked` is not `"true"`. `relabelBlock()` recomputes the effective links and
labels client-side after a toggle or a drag, mirroring `groupSessionExercises` in the controller, so a card dragged to
the bottom of its block loses its chain and rests again. A global lb/kg toggle converts all exercises at once; each non-time-based exercise card also has its own
per-exercise lb/kg toggle that persists the preference to the exercise library via
`POST /sessions/:id/exercises/:eid/unit`. Each exercise card header has a three-dots (⋮) dropdown button — icon shows ↑/↓ for HitMax/BelowGoal — whose menu
contains "Edit goal" (opens the appropriate goal modal) and "Change exercise" (modal with autocomplete that POSTs
`name` to `/sessions/:id/exercises/:eid/change` and reloads the page). Exercise autocomplete fills type radio from `.ExerciseLibraryJSON`.
