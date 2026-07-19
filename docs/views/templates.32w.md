---
level: 32w
parent: ../views.32w.md
deeper: templates-index.128w.md
relates-to:
  - ../controllers/templates.128w.md
source: views/templates/
---

| Template | Summary                                                                                                                                                                                               | Detail                          |
|----------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------|
| index    | Lists workout templates with name and focus; delete-confirmation modal; `.Success` flash                                                                                                              | [128w](templates-index.128w.md) |
| new      | Create-template form: thin shell around `partials/template_form.tpl` (name, focus, dynamic exercise list with drag-to-reorder, type radios, block selector, autocomplete); POSTs to `/templates/new`  | [128w](templates-new.128w.md)   |
| show     | Read-only template detail: exercises grouped by block; "Edit" button; `.Success` flash                                                                                                                | [128w](templates-show.128w.md)  |
| edit     | Edit-template form: the same `partials/template_form.tpl` as new, with edit-page chrome from `c.Data`; POSTs to `/templates/:id`; pre-populated from `.Template` and `.Exercises`                     | [128w](templates-edit.128w.md)  |
