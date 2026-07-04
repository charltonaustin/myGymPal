---
level: 128w
parent: partials.32w.md
deeper: partials-modal_goal_reps.256w.md
relates-to:
  - ../controllers/exercises.128w.md
source: views/partials/modal_goal_reps.tpl
---

Bootstrap modal `#goalRepsModal` for updating the goal rep range of a bodyweight exercise. Triggered by buttons with
`data-ex-name`, `data-goal-rep-min`, and `data-goal-rep-max`. On open, JS populates the exercise name and pre-fills
min/max rep inputs. Save button POSTs AJAX to `/exercises/goal-reps` with `name`, `goal_rep_min`, `goal_rep_max`. On
success, replaces the `X–Y reps` pattern in the card's `.text-muted.small` div and updates the button's data attributes.
