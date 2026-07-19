---
level: 128w
parent: templates.32w.md
deeper: templates-show.256w.md
relates-to:
  - ../controllers/templates.128w.md
source: views/templates/show.tpl
---

Read-only template detail page. Includes `partials/navbar.tpl`. Shows template name and optional focus text. An "Edit"
button links to `/templates/:id/edit`. Below the header, loose exercises are grouped into blocks (`.ExerciseBlocks`);
each block renders a section heading and a list of exercise names (with "Bodyweight" note where applicable).

Each circuit (`.Circuits`) then renders as one card — a dark header with the circuit name and
`N round(s) · Ns transition`, above a list of its exercises with their work durations as badges — rather than as a
loose run of exercises. Circuit members are excluded from `.ExerciseBlocks` by the controller, so an exercise appears
in exactly one place.

Empty state shown only when there are neither blocks nor circuits. `.Success` flash auto-dismissed after 3 s. No AJAX.
