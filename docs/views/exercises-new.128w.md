---
level: 128w
parent: exercises.32w.md
deeper: exercises-new.256w.md
relates-to:
  - ../controllers/exercises.128w.md
source: views/exercises/new.tpl
---

Create-exercise form page. Includes `partials/navbar.tpl`. The form body is entirely rendered via the
`partials/exercise_fields.tpl` partial (name, type radio, conditional goal fields). Error alert auto-removes after 4 s.
Bootstrap `was-validated` is added on submit via JS guard. POSTs to `/exercises/new`.

Template variables: `.Error` (string), and all fields forwarded to `exercise_fields` partial: `.Name`, `.IsBodyweight`,
`.IsTimeBased`, `.GoalWeight`, `.ExWeightUnit`, `.GoalRepMin`, `.GoalRepMax`, `.GoalHours`, `.GoalMinutes`,
`.GoalSecsRemainder`, `.DefaultBlock`, `.ShowDefaultBlock`.
