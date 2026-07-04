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
on each set logged, uses Web Audio API for alarm, and persists across page reloads via `localStorage`. A lb/kg toggle
converts all displayed weights client-side. Exercise autocomplete fills type radio from `.ExerciseLibraryJSON`.
