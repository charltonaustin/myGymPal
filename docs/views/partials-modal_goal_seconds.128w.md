---
level: 128w
parent: partials.32w.md
deeper: partials-modal_goal_seconds.256w.md
relates-to:
  - ../controllers/exercises.128w.md
source: views/partials/modal_goal_seconds.tpl
---

Bootstrap modal `#goalSecondsModal` for updating the goal duration of a time-based exercise. Triggered by buttons with
`data-ex-name` and `data-goal-seconds`. On open, JS converts total seconds to h/m/s and pre-fills three number inputs.
Save POSTs AJAX to `/exercises/goal-seconds` with `name`, `goal_h`, `goal_m`, `goal_s`. On success, formats the returned
`goal_seconds` as a human-readable duration string and updates the card's `.text-muted.small` goal text, then updates
the button's `data-goal-seconds`.
