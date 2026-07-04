---
level: 256w
parent: weight.32w.md
relates-to:
  - ../controllers/weight.128w.md
source: views/weight/index.tpl
---

## Purpose

Log and review body-weight entries with a rolling average summary.

## Partials Included

- `partials/navbar.tpl`

## Template Variables

| Variable       | Type             | Description                                                                         |
|----------------|------------------|-------------------------------------------------------------------------------------|
| `.WeightAvg`   | WeightAvg struct | `.Days` (int), `.Weight` (float64), `.Unit` (string: lb/kg); nil when no entries    |
| `.Entries`     | `[]WeightEntry`  | Each: `.ID` (int), `.Weight` (float64), `.WeightUnit` (string), `.Date` (time.Time) |
| `.DefaultDate` | string           | Pre-fills log-weight date input with today (YYYY-MM-DD)                             |
| `.WeightUnit`  | string           | User's preferred unit; pre-selects the matching `<option>` in the log form          |
| `.Success`     | string           | Flash message; auto-dismissed after 3 000 ms                                        |

## Forms

1. Log Weight form: `method="POST" action="/weight"` — fields: `date` (date, required), `weight` (number, min=0,
   step=0.1, required), `weight_unit` (select: lb/kg, pre-selected by `.WeightUnit`).
2. Per-entry inline Edit form: `method="POST" action="/weight/:id"` — fields: `weight` (number, pre-filled),
   `weight_unit` (select, pre-selected); initially `d-none`.
3. Per-entry Delete form: `method="POST" action="/weight/:id/delete"` — submit button only.

## Conditional Rendering

- `{{if .WeightAvg}}` — renders average card or nothing.
- `{{if .Entries}}` — renders list or empty-state text.
- `{{if eq .WeightUnit "lb"}}` / `{{if eq .WeightUnit "kg"}}` — selects matching option in both the log and edit forms.

## JavaScript Behavior

- Success alert auto-dismiss.
- `showEdit(id)` / `hideEdit(id)`: toggles `d-none` on `.view-row-{id}` and `.edit-row-{id}` within the same `<li>`.

## AJAX / Fetch

None. All mutations use standard POST forms.
