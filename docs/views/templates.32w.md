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
| new      | Create-template form: name, focus, dynamic exercise list with drag-to-reorder, type radios (weighted/bodyweight/time-based), block selector, exercise-library autocomplete; POSTs to `/templates/new` | [128w](templates-new.128w.md)   |
| show     | Read-only template detail: exercises grouped by block; "Edit" button; `.Success` flash                                                                                                                | [128w](templates-show.128w.md)  |
| edit     | Edit-template form identical to new; POSTs to `/templates/:id`; pre-populated from `.Template` and `.Exercises`                                                                                       | [128w](templates-edit.128w.md)  |
