---
level: 256w
parent: templates.32w.md
relates-to:
  - ../controllers/templates.128w.md
source: views/templates/show.tpl
---

## Purpose

Display a workout template's details in read-only form: loose exercises organised by block, and each circuit as its
own grouped card.

## Partials Included

- `partials/navbar.tpl`

## Template Variables

| Variable          | Type                    | Description                                                                                        |
|-------------------|-------------------------|----------------------------------------------------------------------------------------------------|
| `.Template`       | Template struct         | `.ID`, `.Name`, `.Focus`                                                                           |
| `.ExerciseBlocks` | `[]templateExerciseBlock` | Each: `.Label` (display heading), `.Exercises` (`{Name, IsBodyweight}`). **Loose exercises only.** |
| `.Circuits`       | `[]templateCircuitView` | Each: `.Name`, `.Rounds`, `.TransitionSeconds`, `.Exercises` (`{Name, WorkSeconds}`)                |
| `.Success`        | string                  | Flash success message; auto-dismissed after 3 000 ms.                                              |

A circuit's exercises are held back from `.ExerciseBlocks` by `groupTemplateExercises`, which skips any exercise with
a non-nil `CircuitID`. Without that, every circuit exercise would be listed twice — once in its circuit card and again
as a loose exercise in its block.

## Conditional Rendering

- `{{if .Template.Focus}}` — renders focus subtitle.
- `{{range .ExerciseBlocks}}` — block sections for loose exercises.
- `{{range .Circuits}}` — one card per circuit: header shows `{{.Rounds}} round{{if ne .Rounds 1}}s{{end}} ·
  {{.TransitionSeconds}}s transition`; each exercise row carries its `{{.WorkSeconds}}s` badge.
- `{{if and (not .ExerciseBlocks) (not .Circuits)}}` — empty-state paragraph, shown only when there is neither.
- `{{if .IsBodyweight}}` — renders "Bodyweight" secondary text under a loose exercise name.

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
