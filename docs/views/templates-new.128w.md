---
level: 128w
parent: templates.32w.md
deeper: templates-new.256w.md
relates-to:
  - ../controllers/templates.128w.md
source: views/templates/new.tpl
---

Create-template form. Includes `partials/navbar.tpl`. Fields: template name (required), focus (optional). Below that, a
dynamically managed exercise list allows adding/removing rows; each row has a name input, weighted/bodyweight/time-based
radio group, hidden is_bodyweight/is_time_based fields, and a block selector. Rows are drag-to-reorder via SortableJS.
Exercise-name inputs have autocomplete from `.ExerciseLibraryJSON` (type radio is auto-set). A hidden `exercise_count`
field tracks count. POSTs to `/templates/new`.
