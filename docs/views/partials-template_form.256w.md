---
level: 256w
parent: partials.32w.md
relates-to:
  - ../controllers/template.128w.md
  - templates-new.256w.md
  - templates-edit.256w.md
source: views/partials/template_form.tpl
---

## Purpose

The shared body of the create-template and edit-template pages. `views/templates/new.tpl` and
`views/templates/edit.tpl` differ only in `<title>` and five values passed as `c.Data` keys; everything else — the
exercise-row builder, the radio/hidden sync, and the autocomplete JS — lives here, in one file.

## Partials Included

- `partials/navbar.tpl`

## Template Variables

### Page chrome (differs per page — set by the controller)

| Variable       | New template        | Edit template       |
|----------------|---------------------|---------------------|
| `.BackURL`     | `/templates`        | `/templates/:id`    |
| `.BackLabel`   | `Templates`         | template name       |
| `.Heading`     | `New Workout Template` | `Edit Template`  |
| `.FormAction`  | `/templates/new`    | `/templates/:id`    |
| `.SubmitLabel` | `Create Template`   | `Save Changes`      |

`.BackURL` backs both the breadcrumb and the Cancel link — on both pages they point to the same place. Each key is set
by `newFormChrome()` / `editFormChrome()` in `controllers/template.go`, called from both of that page's render paths
(the initial GET and the validation-error re-render). A missing key renders as an empty string, not an error.

### Form data (shared)

| Variable               | Type                 | Description                                                                                      |
|------------------------|----------------------|--------------------------------------------------------------------------------------------------|
| `.Error`               | string               | Server-side error; auto-removed after 4 000 ms.                                                  |
| `.Name`                | string               | Pre-fills template name input.                                                                   |
| `.Focus`               | string               | Pre-fills focus input.                                                                           |
| `.Exercises`           | `[]exerciseForm`     | Exercise rows; each: `.Name`, `.IsBodyweight`, `.IsTimeBased`, `.Block`.                         |
| `.ExerciseLibraryJSON` | template.JS          | JSON array used for exercise-name autocomplete; each entry: `{name, isBodyweight, isTimeBased}`. |

## Form Fields

Form: `id="template-form" method="POST" action="{{.FormAction}}" novalidate`

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
