---
level: 256w
parent: root.32w.md
relates-to:
  - ../controllers/settings.128w.md
source: views/settings.tpl
---

## Purpose

User preference management and account-deletion screen.

## Partials Included

- `partials/navbar.tpl`

## Template Variables

| Variable      | Type   | Description                                                           |
|---------------|--------|-----------------------------------------------------------------------|
| `.Success`    | string | Flash success (e.g. "Settings saved"); auto-dismissed after 3 000 ms. |
| `.Error`      | string | Server-side error; shown as `div.alert.alert-danger`.                 |
| `.WeightUnit` | string | `"lb"` or `"kg"`; pre-checks the matching radio button.               |

## Forms

1. Settings form: `method="post" action="/settings" data-offline-sync`
    - `weight_unit` — radio buttons: `lb` and `kg`.
    - Submit: "Save Settings".
2. Delete Account modal form: `method="post" action="/account/delete"`
    - Submit only: "Yes, Delete Everything".

## Modal

Bootstrap modal `#deleteAccountModal`. Opened by "Delete Account" button (`btn-outline-danger`). Confirm button submits
the delete form. No JS beyond open/close.

## JavaScript Behavior

- Success alert auto-dismiss via `bootstrap.Alert.getOrCreateInstance`.
- No other JS beyond Bootstrap bundle.

## Conditional Rendering

- `{{if .Success}}` — dismissible success alert.
- `{{if .Error}}` — danger alert (not dismissible).
- `{{if eq .WeightUnit "lb"}}` / `{{if eq .WeightUnit "kg"}}` — `checked` attribute on matching radio.

## AJAX / Fetch

None.
