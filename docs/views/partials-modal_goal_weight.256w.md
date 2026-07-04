---
level: 256w
parent: partials.32w.md
relates-to:
  - ../controllers/exercises.128w.md
source: views/partials/modal_goal_weight.tpl
---

## Purpose

Bootstrap modal that lets the user update the goal weight for a weighted exercise directly from the session view,
without navigating away.

## No Template Variables

This partial takes no data from the template context. It is purely DOM-driven via data attributes on the triggering
button.

## Trigger Button Data Attributes (expected on `.hit-max-btn`)

| Attribute          | Description                                                                               |
|--------------------|-------------------------------------------------------------------------------------------|
| `data-ex-name`     | Exercise name (lowercase); sent to server                                                 |
| `data-goal-weight` | Current goal weight (float); pre-fills input                                              |
| `data-weight-unit` | `"lb"` or `"kg"`; pre-selects unit dropdown                                               |
| `data-direction`   | `"up"` (hit max reps) or `"down"` (logged below goal); changes the modal description text |

## Modal Elements

- `#goalWeightModal` — Bootstrap fade modal.
- `#goalWeightExName` — `<strong>` tag injected with exercise name.
- `#goalWeightDesc` — descriptive sentence; first text node changed per direction.
- `#goalWeightInput` — number input, min=0, step=0.5.
- `#goalWeightUnit` — select: lb / kg.
- `#goalWeightError` — hidden error span.
- `#goalWeightSaveBtn` — triggers AJAX save.

## AJAX Endpoint

`POST /exercises/goal-weight`

Payload (URL-encoded): `name`, `goal_weight`, `weight_unit`

Response JSON: `{goal_weight: float, weight_unit: string}` or `{error: string}`

## DOM Updates on Success

1. Finds `.goal-weight-val` within the triggering card; updates `data-w`, `data-u`, and text content.
2. Falls back to replacing a `Goal: N unit` substring in the nearest `.text-muted.small` if no `.goal-weight-val`
   exists.
3. Updates `activeBtn.dataset.goalWeight` and `activeBtn.dataset.weightUnit` for subsequent modal opens.
4. Swaps button icon from `bi-pencil` to `bi-dash-circle-fill` after a goal is set.
5. Pre-fills the card's `input[name="actual_weight"]` with the rounded new goal weight if currently empty or zero.
6. Sets `select[name="weight_unit"]` to the new unit.
7. Hides the modal.
