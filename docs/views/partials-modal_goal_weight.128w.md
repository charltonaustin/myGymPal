---
level: 128w
parent: partials.32w.md
deeper: partials-modal_goal_weight.256w.md
relates-to:
  - ../controllers/exercises.128w.md
source: views/partials/modal_goal_weight.tpl
---

Bootstrap modal `#goalWeightModal` for updating a weighted exercise's goal weight. Opened from session-show buttons that
carry `data-ex-name`, `data-goal-weight`, `data-weight-unit`, and `data-direction` (up/down) attributes. On open, JS
populates the exercise name and pre-fills the weight input. On save, AJAX POSTs to `/exercises/goal-weight` with `name`,
`goal_weight`, `weight_unit`. On success, updates the card's `.goal-weight-val` span, the button's data attributes, and
pre-fills the weight input in the log form. Supports directional descriptions (hit max / logged below goal).
