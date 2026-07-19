---
level: 128w
parent: templates.32w.md
deeper: templates-edit.256w.md
relates-to:
  - ../controllers/templates.128w.md
  - partials-template_form.128w.md
source: views/templates/edit.tpl
---

Edit-template form. A thin shell: `<head>` (title "Edit {{.Template.Name}} — My Gym Pal") plus
`{{template "partials/template_form.tpl" .}}` — the same partial the create page uses, so the exercise list,
drag-to-reorder, autocomplete and JS helpers are the same code rather than a copy of it. See
[partials-template_form.128w.md](partials-template_form.128w.md).

The controller supplies the keys that make the shared partial render as the *edit* page: `.Heading` "Edit Template",
`.FormAction` and `.BackURL` `/templates/:id`, `.SubmitLabel` "Save Changes", `.BackLabel` the template name. Fields are
pre-populated from `.Name`, `.Focus` and `.Exercises`. POSTs to `/templates/:id`.
