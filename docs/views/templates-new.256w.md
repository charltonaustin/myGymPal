---
level: 256w
parent: templates.32w.md
relates-to:
  - ../controllers/templates.128w.md
  - partials-template_form.256w.md
source: views/templates/new.tpl
---

## Purpose

Create a new workout template with a name, optional focus, and a list of exercises (with type and section block per
exercise).

## Structure

The file is a shell. It contributes only the document scaffolding — `<!DOCTYPE>`, `<head>` with the title
"New Template — My Gym Pal", the Bootstrap/icon links, the manifest link, and the drag-handle style — and then includes
the shared body:

```
{{template "partials/template_form.tpl" .}}
```

`views/templates/edit.tpl` is the same shell with a different `<title>`. The form, the exercise rows, and all the
JavaScript live in the partial and are documented in
[partials-template_form.256w.md](partials-template_form.256w.md) — including the full field table and the JS helper
list, which were previously duplicated between this file and `templates-edit.256w.md`.

## Template Variables

Everything the partial needs (`.Error`, `.Name`, `.Focus`, `.Exercises`, `.ExerciseLibraryJSON`), plus the page-chrome
keys that make it render as the create page rather than the edit page:

| Variable       | Value on this page     |
|----------------|------------------------|
| `.Heading`     | `New Workout Template` |
| `.FormAction`  | `/templates/new`       |
| `.SubmitLabel` | `Create Template`      |
| `.BackURL`     | `/templates`           |
| `.BackLabel`   | `Templates`            |

These are set by `newFormChrome()` in `controllers/template.go`, called from **both** render paths that reach this
template: `New()` (the initial GET) and `Create()`'s `renderForm` (the validation-error re-render). A key set on only
one of them renders as an empty string on the other — a blank submit button, with no error and no failing compile.

## AJAX / Fetch

None.
