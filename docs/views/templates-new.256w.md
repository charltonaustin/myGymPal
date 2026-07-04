---
level: 256w
parent: templates.32w.md
relates-to:
  - ../controllers/templates.128w.md
source: views/templates/new.tpl
---

## Purpose

Create a new workout template with a name, optional focus, and a list of exercises (with type and section block per
exercise).

## Partials Included

- `partials/navbar.tpl`

## Template Variables

| Variable               | Type                 | Description                                                                                      |
|------------------------|----------------------|--------------------------------------------------------------------------------------------------|
| `.Error`               | string               | Server-side error; auto-removed after 4 000 ms.                                                  |
| `.Name`                | string               | Pre-fills template name input (on re-render after error).                                        |
| `.Focus`               | string               | Pre-fills focus input.                                                                           |
| `.Exercises`           | `[]TemplateExercise` | Pre-populated exercise rows; each: `.Name`, `.IsBodyweight`, `.IsTimeBased`, `.Block`.           |
| `.ExerciseLibraryJSON` | template.JS          | JSON array used for exercise-name autocomplete; each entry: `{name, isBodyweight, isTimeBased}`. |

## Form Fields

Form: `id="template-form" method="POST" action="/templates/new" novalidate`

Static fields:
| Field | Name | Type | Notes |
|-------|------|------|-------|
| Template Name | `name` | text | required |
| Focus | `focus` | text | optional |
| Exercise Count | `exercise_count` | hidden | managed by JS; submitted count used to parse exercise array |

Per exercise (index `i`):
| Field | Name | Type | Notes |
|-------|------|------|-------|
| Name | `exercise_name_i` | text | required; has autocomplete |
| Type | `ex_type_i` | radio | values: weighted / bodyweight / time_based |
| Is Bodyweight | `is_bodyweight_i` | hidden | "on" or "" — synced by JS |
| Is Time Based | `is_time_based_i` | hidden | "on" or "" — synced by JS |
| Block | `block_i` | select | main / abs / cardio / stretch |

## JavaScript Behavior

- `syncHiddens(row)`: reads selected radio and sets `is_bodyweight_i` / `is_time_based_i` hidden fields.
- `autofillFromLibrary(row, name)`: looks up exercise library entry and checks matching type radio, then calls
  `syncHiddens`.
- `attachAutocomplete(input, row)`: attaches typeahead dropdown to a name input using `exerciseLibraryArr`.
- `bindRow(row)`: attaches remove button, autocomplete, and radio change listeners.
- `reindexRows()`: renumbers all exercise row field names and `exercise_count` after add/remove/reorder.
- "Add Exercise" button: creates a new row HTML, appends to `#exercises-container`, calls `bindRow`.
- SortableJS on `#exercises-container` with `onEnd: reindexRows`.
- Form submit: calls `reindexRows()` then `checkValidity()` guard; adds `was-validated`.

## AJAX / Fetch

None. Exercise library data is embedded in the page as `template.JS` to prevent HTML escaping.
