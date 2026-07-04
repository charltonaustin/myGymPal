---
level: 256w
parent: programs.32w.md
relates-to:
  - ../controllers/programs.128w.md
source: views/programs/index.tpl
---

## Purpose

Lists all training programs for the logged-in user with navigation to create or delete programs.

## Partials Included

- `partials/navbar.tpl`

## Template Variables

| Variable    | Type        | Description                                                                                                                             |
|-------------|-------------|-----------------------------------------------------------------------------------------------------------------------------------------|
| `.Programs` | `[]Program` | Slice of program structs; each has `.ID` (int), `.Name` (string), `.StartDate` (time.Time), `.NumPhases` (int), `.WeeksPerPhase` (int). |
| `.Success`  | string      | Flash success message (e.g. "Program deleted"). Shown in `div.alert.alert-success.alert-dismissible`; auto-dismissed after 3 000 ms.    |

## Conditional Rendering

- `{{if .Programs}}` — renders list group; else shows empty-state paragraph linking to `/programs/new`.
- `{{if gt .NumPhases 1}}` / `{{if gt .WeeksPerPhase 1}}` — pluralises "phase(s)" and "week(s)".

## User Interactions

- "New Program" button → `/programs/new`.
- Program name row → `/programs/:id`.
- Trash icon → opens `#deleteModal`; `data-program-id` and `data-program-name` passed via data attributes.

## Delete Modal

Bootstrap modal `#deleteModal`. On `show.bs.modal` event, JS reads `relatedTarget.dataset.programName` and `programId`,
sets `#deleteModalName` text and `#deleteForm.action` to `/programs/:id/delete`. Submits via standard POST form (no
AJAX).

## JavaScript Behavior

- Success-alert auto-dismiss with `bootstrap.Alert.getOrCreateInstance`.
- Modal `show.bs.modal` listener patches the form action.

## AJAX / Fetch

None.

## Flash Messages

`.Success` only. No `.Error` flash on this page (errors shown on the form pages).
