---
level: 256w
parent: partials.32w.md
relates-to:
  - ../controllers/exercises.128w.md
source: views/partials/modal_goal_reps.tpl
---

## Purpose

Bootstrap modal to update the goal rep range for a bodyweight exercise from the session view.

## No Template Variables

DOM-driven via data attributes on the triggering button.

## Trigger Button Data Attributes

| Attribute           | Description                                    |
|---------------------|------------------------------------------------|
| `data-ex-name`      | Exercise name; sent to server                  |
| `data-goal-rep-min` | Current min reps; pre-fills `#goalRepMinInput` |
| `data-goal-rep-max` | Current max reps; pre-fills `#goalRepMaxInput` |

## Modal Elements

- `#goalRepsModal` — Bootstrap fade modal.
- `#goalRepsExName` — `<strong>` tag with exercise name.
- `#goalRepMinInput` — number input, min=0, step=1.
- `#goalRepMaxInput` — number input, min=0, step=1.
- `#goalRepsError` — hidden error span.
- `#goalRepsSaveBtn` — triggers AJAX save.

## AJAX Endpoint

`POST /exercises/goal-reps`

Payload (URL-encoded): `name`, `goal_rep_min`, `goal_rep_max`

Response JSON: `{goal_rep_min: int, goal_rep_max: int}` or `{error: string}`

## DOM Updates on Success

1. Finds `.text-muted.small` within the triggering card.
2. Replaces the `\d+–\d+ reps` substring with `{goal_rep_min}–{goal_rep_max} reps`.
3. Updates `activeRepsBtn.dataset.goalRepMin` and `.goalRepMax` for subsequent opens.
4. Hides the modal.

## Error Handling

On `data.error`, shows `#goalRepsError` with the message. On `catch`, shows generic error text.
