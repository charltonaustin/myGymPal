---
level: 256w
parent: templates.32w.md
relates-to:
  - ../controllers/templates.128w.md
  - partials-template_form.256w.md
source: views/templates/edit.tpl
---

## Purpose

Edit an existing workout template's name, focus, and exercise list (with type and block per exercise).

## Structure

The file is a shell: document scaffolding and a `<title>` of "Edit {{.Template.Name}} — My Gym Pal", then

```
{{template "partials/template_form.tpl" .}}
```

This is the *same* partial `views/templates/new.tpl` includes. The two pages are no longer independent copies, so the
per-exercise field schema (`exercise_name_i`, `ex_type_i`, `is_bodyweight_i`, `is_time_based_i`, `block_i`,
`exercise_count`) and the JS (`syncHiddens`, `autofillFromLibrary`, `attachAutocomplete`, `bindRow`, `reindexRows`,
SortableJS) exist in exactly one place. They are documented in
[partials-template_form.256w.md](partials-template_form.256w.md).

## Template Variables

| Variable               | Type             | Description                                                                  |
|------------------------|------------------|------------------------------------------------------------------------------|
| `.Template`            | Template struct  | `.Name` is used by the `<title>` in this shell                               |
| `.Name`                | string           | Pre-fills template name input                                                |
| `.Focus`               | string           | Pre-fills focus input                                                        |
| `.Exercises`           | `[]exerciseForm` | Pre-populated rows; each: `.Name`, `.IsBodyweight`, `.IsTimeBased`, `.Block` |
| `.Error`               | string           | Server-side error; auto-removed after 4 000 ms                               |
| `.ExerciseLibraryJSON` | template.JS      | JSON array for autocomplete                                                  |

Page chrome — the keys that make the shared partial render as the edit page:

| Variable       | Value on this page |
|----------------|--------------------|
| `.Heading`     | `Edit Template`    |
| `.FormAction`  | `/templates/:id`   |
| `.SubmitLabel` | `Save Changes`     |
| `.BackURL`     | `/templates/:id`   |
| `.BackLabel`   | template name      |

Set by `editFormChrome(tmpl)` in `controllers/template.go`, called from **both** paths that reach this template:
`Edit()` (the initial GET) and `Update()`'s `renderForm` (the validation-error re-render). Both take the *stored*
template, not the submitted form, so a rejected rename still shows the saved name in the breadcrumb.

## Differences from New Template

Only the `<title>` and the five chrome keys above. Everything else is shared.

## AJAX / Fetch

None.
