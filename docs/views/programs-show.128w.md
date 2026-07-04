---
level: 128w
parent: programs.32w.md
deeper: programs-show.256w.md
relates-to:
  - ../controllers/programs.128w.md
source: views/programs/show.tpl
---

Program detail page. Includes `partials/navbar.tpl`. Shows program name, start date metadata, and two action buttons: "
Start Workout" → `/programs/:id/sessions/new` and "Log Workout" → same path with `?sequential=1`. Displays workout
history (list of sessions with Phase, Week, date, delete button). Below that, a phase-settings editor renders one row
per phase with rep-min, rep-max, sets, rest minutes, rest seconds. A "Copy to all" button clones one phase's values to
all others via JavaScript. POSTs phase settings to `/programs/:id`.
