---
level: 128w
parent: sessions.32w.md
deeper: sessions-new.256w.md
relates-to:
  - ../controllers/sessions.128w.md
source: views/sessions/new.tpl
---

Renders the session-start form. Includes `partials/navbar.tpl`. Heading switches between "Start Workout" and "Log
Workout" based on `.LogMode`. Fields: phase number, week number, workout number, date, optional template selector (
hidden when no templates). A yellow deload notice appears when the selected week equals `weeksPerPhase` (the last week
of a phase). POSTs to `/programs/:id/sessions`.

Template variables: `.Program` (ID, Name, WeeksPerPhase), `.PhaseNumber`, `.WeekNumber`, `.WorkoutNumber`,
`.DefaultDate`, `.Templates` (slice), `.LogMode` (bool). No AJAX.
