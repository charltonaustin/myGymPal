---
level: 256w
parent: exercises.32w.md
relates-to:
  - ../controllers/exercises.128w.md
source: views/exercises/edit.tpl
---

## Purpose

Edit an existing exercise's name, type, and goal values. Mirrors `exercises/new.tpl` in structure; key differences are
the form action and submit label.

## Partials Included

- `partials/navbar.tpl`
- `partials/exercise_fields.tpl` (inside `div.card.p-3`)

## Template Variables

| Variable                      | Type            | Description                                                                                                                                                                                                         |
|-------------------------------|-----------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `.Exercise`                   | Exercise struct | `.ID` used to construct form action `/exercises/:id/edit`. The partial receives the same struct to pre-populate all fields.                                                                                         |
| `.Error`                      | string          | Server-side error; auto-removed from DOM after 4 000 ms.                                                                                                                                                            |
| All exercise-fields variables | —               | Same as new form: `.Name`, `.IsBodyweight`, `.IsTimeBased`, `.GoalWeight`, `.ExWeightUnit`, `.GoalRepMin`, `.GoalRepMax`, `.GoalHours`, `.GoalMinutes`, `.GoalSecsRemainder`, `.DefaultBlock`, `.ShowDefaultBlock`. |

## Form

Form: `id="exercise-form" method="POST" action="/exercises/:id/edit" novalidate`

All inputs rendered by `exercise_fields.tpl`. Page adds Submit ("Save Changes") and Cancel (link to `/exercises`)
buttons.

## JavaScript Behavior

Identical to `exercises/new.tpl`:

- Error alert removed after 4 000 ms.
- Form submit: `checkValidity()` guard; adds `was-validated`.
- `exercise_fields.tpl` provides its own script for type-radio toggling.

## AJAX / Fetch

None.
