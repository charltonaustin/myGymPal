---
level: 128w
parent: root.32w.md
deeper: root-settings.256w.md
relates-to:
  - ../controllers/settings.128w.md
source: views/settings.tpl
---

Account Settings page. Includes `partials/navbar.tpl`. Shows `.Success` (auto-dismissed) and `.Error` alerts. A form
POSTs `weight_unit` (lb/kg radio buttons) to `/settings`. Below a divider, a "Danger Zone" section has a "Delete
Account" button that opens a Bootstrap confirmation modal; the modal form POSTs to `/account/delete`. No AJAX. Has
`data-offline-sync` attribute on the settings form.

Template variables: `.Success`, `.Error`, `.WeightUnit` (string, pre-selects the matching radio).
