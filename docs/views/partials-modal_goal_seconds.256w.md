---
level: 256w
parent: partials.32w.md
relates-to:
  - ../controllers/exercises.128w.md
source: views/partials/modal_goal_seconds.tpl
---

## Purpose

Bootstrap modal to update the goal duration for a time-based exercise from the session view.

## No Template Variables

DOM-driven via data attributes on the triggering button.

## Trigger Button Data Attributes

| Attribute           | Description                                                    |
|---------------------|----------------------------------------------------------------|
| `data-ex-name`      | Exercise name; sent to server                                  |
| `data-goal-seconds` | Current goal as total seconds; converted to h/m/s for pre-fill |

## Modal Elements

- `#goalSecondsModal` — Bootstrap fade modal.
- `#goalSecondsExName` — `<strong>` tag with exercise name.
- `#goalSecondsH` — number input hours (min=0).
- `#goalSecondsM` — number input minutes (min=0, max=59).
- `#goalSecondsS` — number input seconds (min=0, max=59).
- `#goalSecondsError` — hidden error span.
- `#goalSecondsSaveBtn` — triggers AJAX save.

## Helper Functions

- `secsToHMS(total)`: returns `{h, m, s}` from total seconds.
- `fmtHMS(h, m, s)`: formats as `"Xh MMm SSs"`, `"Mm SSs"`, or `"Ss"` depending on magnitude.

## AJAX Endpoint

`POST /exercises/goal-seconds`

Payload (URL-encoded): `name`, `goal_h`, `goal_m`, `goal_s`

Response JSON: `{goal_seconds: int}` or `{error: string}`

## DOM Updates on Success

1. Finds `.text-muted.small` within the triggering button's card.
2. If it contains `Goal:`, replaces the `Goal: ...` portion with `Goal: {fmtHMS(...)}`.
3. Otherwise sets the element's full text content to the new goal label.
4. Updates `activeSecsBtn.dataset.goalSeconds` with the returned `goal_seconds`.
5. Hides the modal.

## Error Handling

On `data.error`, shows `#goalSecondsError`. On network `catch`, shows generic message.
