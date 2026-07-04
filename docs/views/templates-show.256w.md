---
level: 256w
parent: templates.32w.md
relates-to:
  - ../controllers/templates.128w.md
source: views/templates/show.tpl
---

## Purpose

Display a workout template's details in read-only form, organised by exercise block.

## Partials Included

- `partials/navbar.tpl`

## Template Variables

| Variable          | Type              | Description                                                                      |
|-------------------|-------------------|----------------------------------------------------------------------------------|
| `.Template`       | Template struct   | `.ID`, `.Name`, `.Focus`                                                         |
| `.ExerciseBlocks` | `[]ExerciseBlock` | Each: `.Label` (display heading), `.Exercises` (slice of `{Name, IsBodyweight}`) |
| `.Success`        | string            | Flash success message; auto-dismissed after 3 000 ms.                            |

## Conditional Rendering

- `{{if .Template.Focus}}` — renders focus subtitle.
- `{{if .ExerciseBlocks}}` — renders block sections or empty-state paragraph.
- `{{if .IsBodyweight}}` — renders "Bodyweight" secondary text under exercise name.

## User Interactions

- "Edit" button → `/templates/:id/edit`.
- Back link → `/templates`.

## JavaScript Behavior

Success alert auto-dismiss only. No other JS beyond Bootstrap bundle.

## AJAX / Fetch

None.

## Notes

Template exercises on this page show only name and bodyweight flag — no goal weights or reps. Those values come from the
Exercise library and Phase config at session-render time (see CLAUDE.md Key Invariants).
