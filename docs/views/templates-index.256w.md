---
level: 256w
parent: templates.32w.md
relates-to:
  - ../controllers/templates.128w.md
source: views/templates/index.tpl
---

## Purpose

Lists all workout templates for the user; entry point to create, view, or delete templates.

## Partials Included

- `partials/navbar.tpl`

## Template Variables

| Variable     | Type         | Description                                                                                         |
|--------------|--------------|-----------------------------------------------------------------------------------------------------|
| `.Templates` | `[]Template` | Each: `.ID` (int), `.Name` (string), `.Focus` (string, optional)                                    |
| `.Success`   | string       | Flash success; shown in `div.alert.alert-success.alert-dismissible`; auto-dismissed after 3 000 ms. |

## Conditional Rendering

- `{{if .Templates}}` — list group or empty-state text.
- `{{if .Focus}}` — secondary line under template name.

## User Interactions

- "+ New Template" button → `/templates/new`.
- Template name row → `/templates/:id`.
- Trash icon → opens `#deleteModal` with `data-template-id` and `data-template-name`.

## Delete Modal

Bootstrap modal `#deleteModal`. On `show.bs.modal`: sets `#deleteModalName` inner text; sets `#deleteForm.action` to
`/templates/:id/delete`. Standard POST (no AJAX).

## JavaScript Behavior

- Success alert auto-dismiss via `bootstrap.Alert.getOrCreateInstance`.
- Modal `show.bs.modal` listener patches form action.

## AJAX / Fetch

None.
