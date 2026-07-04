---
level: 256w
parent: root.32w.md
relates-to:
  - ../controllers/dashboard.128w.md
source: views/dashboard.tpl
---

## Purpose

Post-login home screen summarising recent health and training activity.

## Partials Included

- `partials/navbar.tpl`

## Template Variables

| Variable          | Type                | Description                                                                                                                   |
|-------------------|---------------------|-------------------------------------------------------------------------------------------------------------------------------|
| `.Username`       | string              | Shown in "Welcome back, {Username}!" heading                                                                                  |
| `.WeightAvg`      | WeightAvg struct    | `.Weight` (float64), `.Unit` (string), `.Days` (int); nil when no entries                                                     |
| `.MacroSummary`   | MacroSummary struct | `.Days`, `.HasGoal`, `.Protein.Actual`, `.Protein.Pct`, `.Protein.AtGoal`; same for Carbs, Fat, Calories; nil when no entries |
| `.RecentSessions` | `[]SessionSummary`  | Each: `.ID`, `.ProgramID`, `.ProgramName`, `.PhaseNumber`, `.WeekNumber`, `.WorkoutNumber`, `.IsDeload`, `.Date`              |

## Conditional Rendering

- `{{if .WeightAvg}}` — shows weight average card content or "No entries yet."
- `{{if .MacroSummary}}` — shows macro table or "No entries yet."
- `{{if .MacroSummary.HasGoal}}` — adds a third column with percentage (green when `AtGoal`, red otherwise).
- `{{if .RecentSessions}}` — renders session list and "Log next workout" link, or empty-state text.
- `{{if .IsDeload}}` — appends "· Deload" to session metadata row.

## User Interactions

- "Log Weight" link → `/weight`.
- "Log Macros" link → `/macros`.
- "Log next workout" link → `/programs/{first-session.ProgramID}/sessions/new`.
- Session rows → `/sessions/:id`.
- Empty-state "Start a session" link → `/programs`.

## JavaScript Behavior

Loads `offline-sync.js` and registers `/sw.js`. No other JS.

## AJAX / Fetch

None.

## Flash Messages

None.
