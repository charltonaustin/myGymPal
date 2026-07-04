---
level: 128w
parent: partials.32w.md
deeper: partials-exercise_fields.256w.md
relates-to:
  - ../controllers/exercises.128w.md
source: views/partials/exercise_fields.tpl
---

Reusable form fragment for exercise creation and editing. Renders: name input (`#ex_name`, `name="name"`), a three-way
type radio group (weighted/bodyweight/time-based), conditional goal rows (`.ex-weight-row` with weight+unit,
`.ex-bw-row` with min–max reps, `.ex-time-row` with h:m:s), and an optional default-section selector (
`name="default_block"`) guarded by `.ShowDefaultBlock`. Hidden fields `is_bodyweight` and `is_time_based` are synced by
inline JS that toggles row visibility on radio change. Variables: `.Name`, `.IsBodyweight`, `.IsTimeBased`,
`.GoalWeight`, `.ExWeightUnit`, `.GoalRepMin`, `.GoalRepMax`, `.GoalHours`, `.GoalMinutes`, `.GoalSecsRemainder`,
`.DefaultBlock`, `.ShowDefaultBlock`.
