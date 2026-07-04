---
level: 128w
parent: exercises.32w.md
deeper: exercises-index.256w.md
relates-to:
  - ../controllers/exercises.128w.md
source: views/exercises/index.tpl
---

Exercise Library list page. Includes `partials/navbar.tpl`. Displays each exercise with name (capitalized), type
badges (Bodyweight, Time-based), and goal (weight+unit, rep range, or duration). A pencil icon links to
`/exercises/:id/edit`. A trash icon opens a Bootstrap delete-confirmation modal that sets a hidden form's action to
`/exercises/:id/delete` and submits via POST.

Template variables: `.Exercises` (slice with `.ID`, `.Name`, `.IsBodyweight`, `.IsTimeBased`, `.GoalWeight`,
`.WeightUnit`, `.GoalRepMin`, `.GoalRepMax`, `.GoalSeconds`), `.Success` (flash, auto-dismissed after 3 s). Uses
`fmtDuration` template function.
