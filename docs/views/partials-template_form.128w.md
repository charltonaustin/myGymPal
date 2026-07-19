---
level: 128w
parent: partials.32w.md
deeper: partials-template_form.256w.md
relates-to:
  - ../controllers/template.128w.md
  - templates-new.128w.md
  - templates-edit.128w.md
source: views/partials/template_form.tpl
---

The entire body of both template forms: navbar, breadcrumb, heading, error alert, name/focus inputs, the dynamic
exercise list (name input with library autocomplete, weighted/bodyweight/time-based radios, hidden
`is_bodyweight_i`/`is_time_based_i` fields, block selector), drag-to-reorder via SortableJS, and all the JS
(`syncHiddens`, `autofillFromLibrary`, `attachAutocomplete`, `bindRow`, `reindexRows`).

`views/templates/new.tpl` and `views/templates/edit.tpl` are thin shells that supply `<head>`/`<title>` and include
this. Everything that differs between the two pages arrives as a `c.Data` key — `.BackURL`, `.BackLabel`, `.Heading`,
`.FormAction`, `.SubmitLabel` — set by the controller. A key the controller forgets renders as an empty string.
