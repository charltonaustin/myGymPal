---
level: 256w
parent: sessions.32w.md
relates-to:
  - ../controllers/sessions.128w.md
source: views/sessions/new.tpl
---

## Purpose

Pre-workout configuration screen that determines the session's metadata before creating it in the database.

## Partials Included

- `partials/navbar.tpl`

## Template Variables

| Variable         | Type           | Description                                                                   |
|------------------|----------------|-------------------------------------------------------------------------------|
| `.Program`       | Program struct | `.ID`, `.Name`, `.WeeksPerPhase` (used to compute deload)                     |
| `.PhaseNumber`   | int            | Pre-fills phase input                                                         |
| `.WeekNumber`    | int            | Pre-fills week input                                                          |
| `.WorkoutNumber` | int            | Pre-fills workout-number input                                                |
| `.DefaultDate`   | string         | Pre-fills date input (today's date, YYYY-MM-DD)                               |
| `.Templates`     | `[]Template`   | List of available templates; if empty, the template selector is hidden        |
| `.LogMode`       | bool           | When true: heading and button read "Log Workout"; when false: "Start Workout" |

## Form Fields

Form: `method="POST" action="/programs/:id/sessions"`

| Field     | Name             | Type   | Constraints                                              |
|-----------|------------------|--------|----------------------------------------------------------|
| Phase     | `phase_number`   | number | min=1, required                                          |
| Week      | `week_number`    | number | min=1, max=WeeksPerPhase, required                       |
| Workout # | `workout_number` | number | min=1, required                                          |
| Date      | `date`           | date   | required                                                 |
| Template  | `template_id`    | select | optional; option "No template" (value="") always present |

## JavaScript Behavior

- `updateDeloadNotice()`: compares `parseInt(weekInput.value)` with `weeksPerPhase` (embedded as a Go template literal).
  Shows or hides `#deload-notice` (yellow warning text) accordingly.
- Fires on `weekInput` `input` event and once on page load.

## Conditional Rendering

- `{{if .Templates}}` — renders template selector only when templates exist.
- `{{if .LogMode}}` — toggles heading and submit button label.

## AJAX / Fetch

None.
