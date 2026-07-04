---
level: 256w
parent: partials.32w.md
relates-to:
  - ../controllers/exercises.128w.md
source: views/partials/exercise_fields.tpl
---

## Purpose

Reusable HTML fragment containing all input fields needed to create or edit an exercise. Included in
`exercises/new.tpl`, `exercises/edit.tpl`, and inside the "Add Exercise" card on `sessions/show.tpl`.

## Template Variables

| Variable             | Type    | Description                                                                            |
|----------------------|---------|----------------------------------------------------------------------------------------|
| `.Name`              | string  | Pre-fills exercise name input                                                          |
| `.IsBodyweight`      | bool    | Checks the Bodyweight radio; shows `.ex-bw-row`; sets `is_bodyweight` hidden to "on"   |
| `.IsTimeBased`       | bool    | Checks the Time-based radio; shows `.ex-time-row`; sets `is_time_based` hidden to "on" |
| `.GoalWeight`        | float64 | Pre-fills goal weight input (shown for weighted type)                                  |
| `.ExWeightUnit`      | string  | Pre-selects lb/kg in weight-unit select                                                |
| `.GoalRepMin`        | int     | Pre-fills rep-range min input                                                          |
| `.GoalRepMax`        | int     | Pre-fills rep-range max input                                                          |
| `.GoalHours`         | int     | Pre-fills hours field of duration goal                                                 |
| `.GoalMinutes`       | int     | Pre-fills minutes field                                                                |
| `.GoalSecsRemainder` | int     | Pre-fills seconds field                                                                |
| `.DefaultBlock`      | string  | Pre-selects section (main/abs/cardio/stretch)                                          |
| `.ShowDefaultBlock`  | bool    | When true, renders the Default Section selector                                        |

## Form Fields

| Field             | Name                               | Type   | Visibility                                              |
|-------------------|------------------------------------|--------|---------------------------------------------------------|
| Exercise Name     | `name`                             | text   | Always; `id="ex_name"` for autocomplete hook            |
| Type (weighted)   | `ex_type_radio` value=`weighted`   | radio  | Always                                                  |
| Type (bodyweight) | `ex_type_radio` value=`bodyweight` | radio  | Always                                                  |
| Type (time-based) | `ex_type_radio` value=`time_based` | radio  | Always                                                  |
| Is Bodyweight     | `is_bodyweight`                    | hidden | `id="ex_is_bw"`; value "on" or ""                       |
| Is Time Based     | `is_time_based`                    | hidden | `id="ex_is_tb"`; value "on" or ""                       |
| Goal Weight       | `goal_weight`                      | number | `.ex-weight-row` (hidden when bodyweight or time-based) |
| Weight Unit       | `weight_unit`                      | select | lb / kg; same row as goal weight                        |
| Goal Rep Min      | `goal_rep_min`                     | number | `.ex-bw-row` (shown only for bodyweight)                |
| Goal Rep Max      | `goal_rep_max`                     | number | `.ex-bw-row`                                            |
| Goal Hours        | `goal_h`                           | number | `.ex-time-row` (shown only for time-based)              |
| Goal Minutes      | `goal_m`                           | number | `.ex-time-row`                                          |
| Goal Seconds      | `goal_s`                           | number | `.ex-time-row`                                          |
| Default Section   | `default_block`                    | select | Optional; shown when `.ShowDefaultBlock` is true        |

## JavaScript Behavior (inline IIFE)

`updateRows()` reads the currently checked `input[name="ex_type_radio"]`:

- Sets `is_bodyweight` to "on" when bodyweight, "" otherwise.
- Sets `is_time_based` to "on" when time_based, "" otherwise.
- Toggles `.d-none` on `.ex-weight-row`, `.ex-time-row`, `.ex-bw-row` to show only the relevant goal section.

Fires immediately and on each radio `change` event.

## Notes

Used in the `sessions/show.tpl` "Add Exercise" card as well, where `.ShowDefaultBlock` is false (block is chosen via a
separate `<select name="block">` outside this partial).
