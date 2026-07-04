---
level: 256w
parent: macros.32w.md
relates-to:
  - ../controllers/macros.128w.md
source: views/macros/index.tpl
---

## Purpose

All-in-one macro tracking page: summary statistics, goal management, food logging, and a scrollable history with inline
editing.

## Partials Included

- `partials/navbar.tpl`

## Template Variables

| Variable            | Type                | Description                                                                                                                            |
|---------------------|---------------------|----------------------------------------------------------------------------------------------------------------------------------------|
| `.Summary`          | MacroSummary struct | `.Days` (int), `.HasGoal` (bool), `.Protein.Actual`, `.Protein.Goal`, `.Protein.Pct`, `.Protein.AtGoal`; same for Carbs, Fat, Calories |
| `.Goal`             | MacroGoal struct    | `.Protein`, `.Carbs`, `.Fat` (all float64); nil when no goal set                                                                       |
| `.Days`             | `[]DayGroup`        | Grouped by date; each: `.Date`, `.Protein`, `.Carbs`, `.Fat`, `.Calories`, `.Entries`                                                  |
| `.Days[].Entries[]` | MacroEntry          | `.ID`, `.FoodName`, `.ServingWeight`, `.ServingUnit`, `.Protein`, `.Carbs`, `.Fat`                                                     |
| `.DefaultDate`      | string              | Pre-fills food-log date input with today (YYYY-MM-DD)                                                                                  |
| `.FoodHistoryJSON`  | template.JS         | JSON array of past food entries for autocomplete; each: `{name, servingWeight, servingUnit, protein, carbs, fat}`                      |

## Forms

1. Daily Goals form: `method="POST" action="/macros/goals"` — fields: `protein_goal`, `carbs_goal`, `fat_goal` (all
   number, step=1).
2. Log Food form: `method="POST" action="/macros"` — fields: `food_name` (text, required), `date` (date, required),
   `serving_weight` (number), `serving_unit` (select: g/oz/ml/fl oz), `protein`, `carbs`, `fat` (number, step=0.1).
3. Per-entry inline Edit form: `method="POST" action="/macros/:id"` — same fields as log form minus date; initially
   `d-none`, shown by `showEdit(id)`.
4. Per-entry Delete form: `method="POST" action="/macros/:id/delete"` — submit button only.

## Conditional Rendering

- `{{if .Summary}}` — renders N-day average summary card.
- `{{if .Summary.HasGoal}}` — adds Goal and % columns to summary table; colours `AtGoal` cells green/red.
- `{{with $.Goal}}` — renders per-day macro totals against goal; else shows bare totals.
- `{{if .Days}}` — renders per-day entry list or empty-state text.
- `{{if gt .ServingWeight 0.0}}` — shows serving weight/unit in entry view.

## JavaScript Behavior

- `showEdit(id)` / `hideEdit(id)`: toggles `d-none` on `.view-row-{id}` and `.edit-row-{id}`.
- Food-name autocomplete: reads `.FoodHistoryJSON`, filters on input, shows dropdown, calls `fillForm(food)` which
  pre-fills serving weight/unit and macro fields from the matching history entry.
- Autocomplete hides on blur (150 ms delay) to allow click on dropdown item.

## AJAX / Fetch

None — all forms use standard POST. Inline edit/hide is purely DOM manipulation.
