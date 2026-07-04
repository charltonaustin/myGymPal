---
level: 256w
parent: templates.32w.md
relates-to:
  - ../controllers/templates.128w.md
source: views/templates/edit.tpl
---

## Purpose

Edit an existing workout template's name, focus, and exercise list (with type and block per exercise).

## Partials Included

- `partials/navbar.tpl`

## Template Variables

| Variable               | Type                 | Description                                                                  |
|------------------------|----------------------|------------------------------------------------------------------------------|
| `.Template`            | Template struct      | `.ID` (used in form action and cancel link), `.Name`                         |
| `.Name`                | string               | Pre-fills template name input                                                |
| `.Focus`               | string               | Pre-fills focus input                                                        |
| `.Exercises`           | `[]TemplateExercise` | Pre-populated rows; each: `.Name`, `.IsBodyweight`, `.IsTimeBased`, `.Block` |
| `.Error`               | string               | Server-side error; auto-removed after 4 000 ms                               |
| `.ExerciseLibraryJSON` | template.JS          | JSON array for autocomplete; same structure as new template page             |

## Form Fields

Form: `id="template-form" method="POST" action="/templates/:id" novalidate`

Identical to `templates/new.tpl` per-exercise field schema (`exercise_name_i`, `ex_type_i`, `is_bodyweight_i`,
`is_time_based_i`, `block_i`) plus `exercise_count` hidden. See `templates-new.256w.md` for full field table.

## JavaScript Behavior

Identical to `templates/new.tpl`:

- `syncHiddens(row)`, `autofillFromLibrary(row, name)`, `attachAutocomplete(input, row)`, `bindRow(row)`,
  `reindexRows()`.
- "Add Exercise" button appends new row.
- SortableJS with `onEnd: reindexRows`.
- Form submit: `reindexRows()`, `checkValidity()` guard, `was-validated`.

## Differences from New Template

| Aspect         | New                     | Edit                          |
|----------------|-------------------------|-------------------------------|
| Form action    | `/templates/new`        | `/templates/:id`              |
| Submit label   | "Create Template"       | "Save Changes"                |
| Cancel target  | `/templates`            | `/templates/:id`              |
| Pre-population | Empty / error re-render | Loaded from existing template |

## AJAX / Fetch

None.
