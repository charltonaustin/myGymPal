---
level: 256w
parent: partials.32w.md
relates-to:
  - ../controllers/template.128w.md
  - templates-new.256w.md
  - templates-edit.256w.md
source: views/partials/template_form.tpl, views/partials/template_exercise_row.tpl
---

## Purpose

The shared body of the create-template and edit-template pages. `views/templates/new.tpl` and
`views/templates/edit.tpl` differ only in `<title>` and five values passed as `c.Data` keys; everything else — the
exercise-row builder, the circuits section, the radio/hidden sync, and the autocomplete JS — lives here, in one file.

## Partials Included

- `partials/navbar.tpl`
- `partials/template_exercise_row.tpl` — one exercise row. Rendered from two places (loose, and nested inside a
  circuit card), so it is a partial rather than markup repeated twice.

## Template Variables

### Page chrome (differs per page — set by the controller)

| Variable       | New template           | Edit template     |
|----------------|------------------------|-------------------|
| `.BackURL`     | `/templates`           | `/templates/:id`  |
| `.BackLabel`   | `Templates`            | template name     |
| `.Heading`     | `New Workout Template` | `Edit Template`   |
| `.FormAction`  | `/templates/new`       | `/templates/:id`  |
| `.SubmitLabel` | `Create Template`      | `Save Changes`    |

`.BackURL` backs both the breadcrumb and the Cancel link — on both pages they point to the same place. Each key is set
by `newFormChrome()` / `editFormChrome()` in `controllers/template.go`, called from both of that page's render paths
(the initial GET and the validation-error re-render). A missing key renders as an empty string, not an error.

### Form body (differs per render path — set by `setFormBody`)

All four render paths (`New`, `Create`'s error re-render, `Edit`, `Update`'s error re-render) build a
`templateFormData` struct and pass it to `setFormBody`, which sets every key below. A forgotten field is a compile
error rather than a blank field on the page.

| Variable               | Type                 | Description                                                                                      |
|------------------------|----------------------|--------------------------------------------------------------------------------------------------|
| `.Error`               | string               | Server-side error; auto-removed after 4 000 ms. (Set per path, not by `setFormBody`.)            |
| `.Name`                | string               | Pre-fills template name input.                                                                   |
| `.Focus`               | string               | Pre-fills focus input.                                                                           |
| `.Exercises`           | `[]exerciseForm`     | **Loose exercises only.** Each: `.Index`, `.Name`, `.IsBodyweight`, `.IsTimeBased`, `.Block`, `.WorkSeconds`. |
| `.Circuits`            | `[]circuitForm`      | Each: `.Index`, `.Name`, `.Rounds`, `.TransitionSeconds`, `.Exercises` (its own rows).           |
| `.ExerciseCount`       | int                  | Total rows across loose **and** circuits — what the server loops over on submit.                 |
| `.CircuitCount`        | int                  | Number of circuit cards.                                                                         |
| `.ExerciseLibraryJSON` | template.JS          | JSON array used for exercise-name autocomplete; each entry: `{name, isBodyweight, isTimeBased}`. |

`.Index` is carried on each row rather than taken from the template's `range` index, because circuit rows are rendered
nested inside their circuit card, where a range index would restart at zero and collide with the loose rows. Indices
are contiguous `0..ExerciseCount-1` across the whole form.

## Form Fields

Form: `id="template-form" method="POST" action="{{.FormAction}}" novalidate`

Static fields:
| Field          | Name             | Type   | Notes                                            |
|----------------|------------------|--------|--------------------------------------------------|
| Template Name  | `name`           | text   | required                                         |
| Focus          | `focus`          | text   | optional                                         |
| Exercise Count | `exercise_count` | hidden | managed by JS; total rows, circuit rows included |
| Circuit Count  | `circuit_count`  | hidden | managed by JS                                    |

Per exercise (index `i`):
| Field           | Name               | Type   | Notes                                                        |
|-----------------|--------------------|--------|--------------------------------------------------------------|
| Name            | `exercise_name_i`  | text   | required; has autocomplete                                   |
| Type            | `ex_type_i`        | radio  | values: weighted / bodyweight / time_based                   |
| Is Bodyweight   | `is_bodyweight_i`  | hidden | "on" or "" — synced by JS                                    |
| Is Time Based   | `is_time_based_i`  | hidden | "on" or "" — synced by JS                                    |
| Block           | `block_i`          | select | main / abs / cardio / stretch; hidden by CSS inside a circuit |
| Work Seconds    | `work_seconds_i`   | number | shown by CSS only inside a circuit; zeroed on loose rows      |
| Circuit Index   | `circuit_index_i`  | hidden | circuit position, or `-1` for loose; written by `reindexAll`  |

Per circuit (index `j`):
| Field      | Name                    | Type   | Notes                     |
|------------|-------------------------|--------|---------------------------|
| Name       | `circuit_name_j`        | text   | required; empty ⇒ dropped |
| Rounds     | `circuit_rounds_j`      | number | `min=1`                   |
| Transition | `circuit_transition_j`  | number | `min=0`, seconds          |

## JavaScript Behavior

- `syncHiddens(row)`: reads selected radio and sets `is_bodyweight_i` / `is_time_based_i` hidden fields.
- `autofillFromLibrary(row, name)`: looks up exercise library entry and checks matching type radio, then calls
  `syncHiddens`.
- `attachAutocomplete(input, row)`: attaches typeahead dropdown to a name input using `exerciseLibraryArr`.
- `bindRow(row)` / `bindCircuit(card)`: attach remove buttons, autocomplete, radio listeners, per-circuit
  "Add Exercise to Circuit", and a SortableJS list for the circuit's rows.
- `reindexCircuits()`: renumbers circuit cards and their field names; updates `circuit_count`.
- `reindexRows()`: renumbers every exercise row's field names and `exercise_count`; reads each row's circuit
  membership out of the DOM (`row.closest('.circuit-card')`) and writes it into `circuit_index_i`; zeroes
  `work_seconds_i` on rows that are not in a circuit.
- `reindexAll()`: circuits first, then rows — rows read the circuit index off the card, so the order matters.
  Runs on load, on every add/remove/reorder, and on submit.
- SortableJS on `#exercises-container` and on each `.circuit-exercises`, with `onEnd: reindexAll`.
- Form submit: calls `reindexAll()` then `checkValidity()` guard; adds `was-validated`.

**Renaming, not just re-valuing.** `reindexRows` sets each hidden field's `name` as well as its `value`. Leaving the
name on the index a row was born with makes row *i* submit row *j*'s membership after any renumber, and the field for
row 0 can go missing entirely — which parses server-side as "not in a circuit". The controller tests cannot catch
this: they build form fields in Go with the names already correct, so the bug lives only in the browser's renumbering
between render and submit.

## Circuit membership is structural

An exercise is in a circuit because it sits inside that circuit's card. There is no per-row circuit selector: a select
would let the DOM and the data disagree the moment a row is dragged. The same reasoning makes the work-seconds and
block fields CSS-toggled on the containing card rather than conditionally rendered — one row markup, and a row keeps
working wherever it ends up.

## AJAX / Fetch

None. Exercise library data is embedded in the page as `template.JS` to prevent HTML escaping.
