---
level: 128w
parent: templates.32w.md
deeper: templates-new.256w.md
relates-to:
  - ../controllers/templates.128w.md
  - partials-template_form.128w.md
source: views/templates/new.tpl
---

Create-template form. A thin shell: `<head>` (title "New Template — My Gym Pal") plus
`{{template "partials/template_form.tpl" .}}`, which carries the whole page — navbar, name and focus inputs, the
dynamic exercise list with drag-to-reorder, type radios, block selectors, and exercise-library autocomplete. See
[partials-template_form.128w.md](partials-template_form.128w.md).

The controller supplies the keys that make the shared partial render as the *create* page: `.Heading` "New Workout
Template", `.FormAction` `/templates/new`, `.SubmitLabel` "Create Template", `.BackURL` `/templates`, `.BackLabel`
"Templates". POSTs to `/templates/new`.
