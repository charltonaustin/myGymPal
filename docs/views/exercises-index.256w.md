---
level: 256w
parent: exercises.32w.md
relates-to:
  - ../controllers/exercises.128w.md
source: views/exercises/index.tpl
---

## Purpose

Displays all exercises in the user's library with their type, goal information, and CRUD actions.

## Partials Included

- `partials/navbar.tpl`

## Template Variables

| Variable     | Type         | Description                                                                                                                                       |
|--------------|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------|
| `.Exercises` | `[]Exercise` | Each: `.ID`, `.Name`, `.IsBodyweight`, `.IsTimeBased`, `.GoalWeight` (float64), `.WeightUnit`, `.GoalRepMin`, `.GoalRepMax`, `.GoalSeconds` (int) |
| `.Success`   | string       | Flash message (e.g. after delete); shown in `div.alert.alert-success.alert-dismissible`; auto-dismissed after 3 000 ms.                           |

## Conditional Rendering

- `{{if .Exercises}}` — list group or empty-state paragraph.
- `{{if .IsBodyweight}}` / `{{if .IsTimeBased}}` — badge display.
- Goal row per exercise: `{{if .IsTimeBased}}` shows `fmtDuration .GoalSeconds`; `{{else if .IsBodyweight}}` shows rep
  range (only when both > 0); `{{else}}` shows `GoalWeight WeightUnit`.

## Template Functions

- `fmtDuration` — formats `.GoalSeconds` as `m:ss` or `h:mm:ss`.
- `printf "%.0f"` — formats `GoalWeight` without decimal.

## User Interactions

- "+ New Exercise" button → `/exercises/new`.
- Pencil icon → `/exercises/:id/edit`.
- Trash icon → opens `#deleteModal` with `data-exercise-id` and `data-exercise-name`.

## Delete Modal

Bootstrap modal `#deleteModal`. On `show.bs.modal`: sets `#deleteModalName` text; sets `#deleteForm.action` to
`/exercises/:id/delete`. POST form submit (no AJAX).

## AJAX / Fetch

None.
