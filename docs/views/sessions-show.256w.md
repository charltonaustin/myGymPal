---
level: 256w
parent: sessions.32w.md
relates-to:
  - ../controllers/sessions.128w.md
source: views/sessions/show.tpl
---

## Purpose

The main workout-logging screen. Displays all exercises grouped into blocks, lets the user log sets, reorder exercises,
adjust goal weights/reps/seconds via modals, and tracks rest time between sets.

## Partials Included

- `partials/navbar.tpl`
- `partials/exercise_fields.tpl` (inside "Add Exercise" card)
- `partials/modal_goal_weight.tpl`
- `partials/modal_goal_reps.tpl`
- `partials/modal_goal_seconds.tpl`

## Template Variables

| Variable                        | Type              | Description                                                                                                                                |
|---------------------------------|-------------------|--------------------------------------------------------------------------------------------------------------------------------------------|
| `.Session`                      | Session struct    | `.ID`, `.WorkoutNumber`, `.PhaseNumber`, `.WeekNumber`, `.Date`, `.IsDeload`                                                               |
| `.Program`                      | Program struct    | `.ID`, `.Name`                                                                                                                             |
| `.ExerciseBlocks`               | `[]ExerciseBlock` | Each block: `.Block` (string: main/abs/cardio/stretch), `.Label`, `.Exercises`                                                             |
| `.ExerciseBlocks[].Exercises[]` | ExerciseEntry     | `.Exercise` (name, ID, goal fields, IsTimeBased, IsBodyweight, LinkedToNext), `.Sets`, `.LastSet`, `.HitMax`, `.BelowGoal`, `.GoalRepMin`, `.GoalRepMax`, `.SupersetLinked`, `.SupersetLabel` |
| `.PhaseRestSeconds`             | int               | Rest timer duration; 0 means timer disabled                                                                                                |
| `.PhaseRepMin` / `.PhaseRepMax` | int               | Phase-level rep targets shown as goal hint                                                                                                 |
| `.WeightUnit`                   | string            | `"lb"` or `"kg"`; user's global preference — drives the global toggle default; each exercise renders in its own `exercises.weight_unit`    |
| `.ExerciseLibraryJSON`          | template.JS       | JSON array of all exercises for autocomplete                                                                                               |

## AJAX Endpoints

| Method | Path                                               | Trigger                 | Payload                                                                                              |
|--------|----------------------------------------------------|-------------------------|------------------------------------------------------------------------------------------------------|
| POST   | `/sessions/:id/exercises/:exId/sets`               | Log-set form submit     | `actual_weight`, `weight_unit`, `actual_reps` OR `actual_h`, `actual_m`, `actual_s`, `activity_type` |
| POST   | `/sessions/:id/exercises/:exId/sets/:setId/delete` | Delete-set form submit  | (none)                                                                                               |
| POST   | `/sessions/:id/exercises/reorder`                  | SortableJS `onEnd`      | `ids=comma-separated-exIds`                                                                          |
| POST   | `/sessions/:id/exercises/:eid/unit`                | Per-exercise unit toggle | `weight_unit=lb\|kg`                                                                                |
| POST   | `/sessions/:id/exercises/:eid/change`              | Change exercise modal    | `name`                                                                                               |
| POST   | `/sessions/:id/exercises/:eid/link`                | Chain (superset) button  | `linked=true\|false`                                                                                 |
| POST   | `/account/unit`                                    | Global unit toggle       | `weight_unit=lb\|kg`                                                                                |
| POST   | `/exercises/goal-weight`                           | Goal weight modal save  | `name`, `goal_weight`, `weight_unit`                                                                 |
| POST   | `/exercises/goal-reps`                             | Goal reps modal save    | `name`, `goal_rep_min`, `goal_rep_max`                                                               |
| POST   | `/exercises/goal-seconds`                          | Goal seconds modal save | `name`, `goal_h`, `goal_m`, `goal_s`                                                                 |

All set-log requests send `X-Requested-With: XMLHttpRequest`. Response is JSON `{id, set_number}`.

## Form Actions (non-AJAX)

- `POST /sessions/:id/exercises` — Add Exercise form
- `POST /sessions/:id/exercises/:exId/delete` — Remove exercise from session
- `POST /sessions/:id/exercises/:exId/cardio/:cardioId/delete` — Remove cardio log entry

## JavaScript Behavior

- **Weight unit toggle**: a global lb/kg radio pair calls `applyGlobalToggle()`, iterating all exercise cards. Each
  weighted exercise card also has its own lb/kg radio pair; changing it calls `applyUnitToCard(card, unit)`, which
  converts `.weight-cell`, `.goal-weight-val`, `[data-goal-weight]` buttons, and the set-log form within that card, then
  fires `POST /sessions/:id/exercises/:eid/unit` to persist the preference to the exercise library. Each card carries
  `data-ex-id` (session_exercise ID) and `data-server-unit` (the unit the server rendered in).
- **Set logging AJAX**: intercepts `.log-set-form` submit; on success, appends a new `<tr>` to the sets table and calls
  `window.startRestTimer()` — but only when the card's `data-linked` is not `"true"`, so a superset flows straight into
  the next exercise with no rest.
- **Supersets**: each card carries `data-link-raw` (stored `linked_to_next`) and `data-linked` (the effective link
  computed by the controller). Every card except the last in its block renders a `.chain-btn`; clicking it POSTs to
  `/sessions/:id/exercises/:eid/link` and, on success, calls `relabelBlock()` — which recomputes effective links and
  `A1`/`A2` badges client-side (mirroring `groupSessionExercises`, including the four-member cap), adds or removes chain
  buttons, and never reloads the page, so a running rest timer survives. A rejected toggle (400) leaves the chain off.
  `relabelBlock()` also runs on SortableJS `onEnd`, so a card dragged to the bottom of its block drops its now-inert
  link and rests again.
- **Delete set AJAX**: event delegation on `.delete-set-form`; removes the row and renumbers remaining sets.
- **SortableJS**: creates a Sortable instance per `.sortable-block`; fires reorder fetch on `onEnd`, then `relabelBlock()`.
- **Rest timer**: IIFE managing a fixed-bottom countdown panel. Stores start time and duration in `localStorage` keyed
  by session ID, restoring on reload. Uses Web Audio API (`AudioContext`) for a three-beep alarm and
  `navigator.serviceWorker.ready` for a push notification when rest completes.
- **Exercise autocomplete**: IIFE that reads `.ExerciseLibraryJSON`, filters as the user types in `#ex_name`, and calls
  `autofillFromLibrary()` to set the type radio and goal fields.

## Conditional Rendering

- `{{if ne .Block "main"}}` — renders block heading for non-main blocks.
- `{{if eq .Block "cardio"}}` — cardio exercises render cardio-log entries instead of sets tables.
- `.HitMax` / `.BelowGoal` — controls which icon appears on the three-dots dropdown trigger button: ↑ for HitMax,
  ↓ for BelowGoal, ⋮ otherwise (time-based exercises always show ⋮). Dropdown contains two items: "Edit goal" (opens
  goalWeightModal/goalRepsModal/goalSecondsModal based on exercise type) and "Change exercise" (opens
  `#changeExerciseModal` with autocomplete input that POSTs to `/sessions/:id/exercises/:eid/change` on submit,
  causing a page reload with the renamed exercise).
- `{{if .Exercise.IsTimeBased}}` — switches between time-based (h:m:s) and weight/reps log forms and set tables.
- `{{if .Session.IsDeload}}` — shows "Deload" badge on session heading.
