---
level: 256w
parent: programs.32w.md
relates-to:
  - ../controllers/programs.128w.md
source: views/programs/show.tpl
---

## Purpose

Program detail view combining workout history and per-phase configuration editing.

## Partials Included

- `partials/navbar.tpl`

## Template Variables

| Variable    | Type           | Description                                                                        |
|-------------|----------------|------------------------------------------------------------------------------------|
| `.Program`  | Program struct | `.ID`, `.Name`, `.StartDate`, `.NumPhases`, `.WeeksPerPhase`                       |
| `.Sessions` | `[]Session`    | Each: `.ID`, `.WorkoutNumber`, `.PhaseNumber`, `.WeekNumber`, `.IsDeload`, `.Date` |
| `.Phases`   | `[]Phase`      | Each: `.PhaseNumber`, `.RepMin`, `.RepMax`, `.DefaultSets`, `.RestSeconds`         |
| `.Success`  | string         | Flash success (e.g. "Settings saved"); auto-dismissed after 3 000 ms.              |
| `.Error`    | string         | Inline danger alert.                                                               |

## Template Functions Used

- `restMinutes` / `restSecs` — custom template functions that split `.RestSeconds` into minutes and remainder seconds.
- `gt`, `if` — conditional pluralisation of "phase(s)", "week(s)".

## Conditional Rendering

- `{{if .Sessions}}` — renders session list or empty-state text.
- `{{if .Session.IsDeload}}` — shows a "Deload" badge.
- `{{if gt .PhaseRestSeconds 0}}` — omitted on programs with no rest configured (this is a session variable but the same
  pattern applies here).

## Phase-Settings Form

Form: `method="POST" action="/programs/:id" data-offline-sync`

One row per phase. Field names are indexed: `rep_min_N`, `rep_max_N`, `sets_N`, `rest_m_N`, `rest_s_N` where N is the
phase number.

## Workout History

Each session row has an inline delete form: `method="POST" action="/sessions/:id/delete"` with a JS `confirm()` dialog.

## JavaScript Behavior

- Success alert auto-dismiss.
- `.copy-to-all` button: reads min, max, sets, rest-m, rest-s from the clicked phase row; writes those values to every
  other phase row's matching inputs.

## AJAX / Fetch

None.
