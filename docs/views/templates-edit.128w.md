---
level: 128w
parent: templates.32w.md
deeper: templates-edit.256w.md
relates-to:
  - ../controllers/templates.128w.md
source: views/templates/edit.tpl
---

Edit-template form. Structurally identical to `templates/new.tpl` — same dynamic exercise list, drag-to-reorder (
SortableJS), exercise-library autocomplete, type radios, block selectors, and JS helpers (`syncHiddens`, `reindexRows`,
`bindRow`). Differences: heading reads "Edit Template", submit button reads "Save Changes", cancel link goes to
`/templates/:id`, form POSTs to `/templates/:id`, and all fields are pre-populated from `.Template` and `.Exercises`.
