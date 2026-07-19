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
exercise list, the circuits section, drag-to-reorder via SortableJS, and all the JS (`syncHiddens`,
`autofillFromLibrary`, `attachAutocomplete`, `bindRow`, `bindCircuit`, `reindexCircuits`, `reindexRows`,
`reindexAll`). One exercise row is `partials/template_exercise_row.tpl`, rendered both loose and inside a circuit.

A circuit is a card with a name, round count and transition, holding exercise rows. **Which circuit an exercise
belongs to is decided by the card it physically sits in, not by a field on the row**: `reindexAll` reads membership
out of the DOM and writes it into each row's hidden `circuit_index_i` before every submit. The work-seconds and block
fields are shown/hidden by CSS on the containing card for the same reason — one row markup, no code path per case.

`reindexAll` must rename fields, not just re-value them. Renumbering rows while leaving a hidden field's `name` on its
original index makes row *i* submit row *j*'s value, and that is a bug no controller test can see, because the tests
build their form fields in Go with the names already correct.

`views/templates/new.tpl` and `views/templates/edit.tpl` are thin shells that supply `<head>`/`<title>` and include
this. Everything that differs between the two pages arrives as a `c.Data` key — `.BackURL`, `.BackLabel`, `.Heading`,
`.FormAction`, `.SubmitLabel` — set by the controller. A key the controller forgets renders as an empty string.
