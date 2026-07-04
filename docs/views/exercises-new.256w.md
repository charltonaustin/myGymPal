---
level: 256w
parent: exercises.32w.md
relates-to:
  - ../controllers/exercises.128w.md
source: views/exercises/new.tpl
---

## Purpose

Form to add a new exercise to the user's library. All field rendering is delegated to the `exercise_fields` partial.

## Partials Included

- `partials/navbar.tpl`
- `partials/exercise_fields.tpl` (inside `div.card.p-3`)

## Template Variables

| Variable                      | Type   | Description                                                                                                                                                                                                                                                                       |
|-------------------------------|--------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `.Error`                      | string | Server-side validation error; removed from DOM after 4 000 ms via `alertEl.remove()`.                                                                                                                                                                                             |
| All exercise-fields variables | —      | Forwarded to partial: `.Name`, `.IsBodyweight`, `.IsTimeBased`, `.GoalWeight`, `.ExWeightUnit`, `.GoalRepMin`, `.GoalRepMax`, `.GoalHours`, `.GoalMinutes`, `.GoalSecsRemainder`, `.DefaultBlock`, `.ShowDefaultBlock`. See `partials-exercise_fields.256w.md` for field details. |

## Form

Form: `id="exercise-form" method="POST" action="/exercises/new" novalidate`

Fields are entirely provided by `exercise_fields.tpl`. The wrapping page adds Submit ("Add Exercise") and Cancel ("Back
to Exercises") buttons.

## JavaScript Behavior

- `alertEl.remove()` after 4 000 ms if error alert exists.
- Form submit: `checkValidity()` guard; adds `was-validated` class.

## AJAX / Fetch

None. The `exercise_fields` partial contains its own inline `<script>` for type-radio toggling (show/hide goal field
rows), but no network calls.
