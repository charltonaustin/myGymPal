---
level: 32w
parent: ../views.32w.md
deeper: partials-navbar.128w.md
relates-to:
  - ../controllers/auth.128w.md
source: views/partials/
---

| Template           | Summary                                                                                                                                                                                                                         | Detail                                      |
|--------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------|
| navbar             | Responsive dark Bootstrap navbar with auth-conditional links (Programs, Templates, Exercises, Weight, Macros, Dashboard, Settings, Log Out) and a dev-environment banner                                                        | [128w](partials-navbar.128w.md)             |
| exercise_fields    | Reusable form fragment: name input, type radio group (weighted/bodyweight/time-based), conditional goal fields (weight+unit, rep range, or h:m:s duration), optional default-section selector; inline JS toggles row visibility | [128w](partials-exercise_fields.128w.md)    |
| modal_goal_weight  | Bootstrap modal for updating an exercise's goal weight via `POST /exercises/goal-weight` (AJAX); updates card DOM on success                                                                                                    | [128w](partials-modal_goal_weight.128w.md)  |
| modal_goal_reps    | Bootstrap modal for updating bodyweight exercise goal rep range via `POST /exercises/goal-reps` (AJAX); updates card DOM on success                                                                                             | [128w](partials-modal_goal_reps.128w.md)    |
| modal_goal_seconds | Bootstrap modal for updating time-based exercise goal duration via `POST /exercises/goal-seconds` (AJAX); updates card DOM on success                                                                                           | [128w](partials-modal_goal_seconds.128w.md) |
