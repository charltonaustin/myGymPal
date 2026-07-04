---
level: 128w
parent: root.32w.md
deeper: root-dashboard.256w.md
relates-to:
  - ../controllers/dashboard.128w.md
source: views/dashboard.tpl
---

Authenticated dashboard. Includes `partials/navbar.tpl`. Greets user by `.Username`. Two summary cards side-by-side:
weight (N-day average, link to `/weight`) and macros (N-day average table with optional goal % columns coloured
green/red). Below, a "Recent Workouts" section lists the last N sessions with program name, phase/week/workout numbers,
deload flag, and date. A "Log next workout" link appears beside the heading when sessions exist. No forms. No AJAX.
