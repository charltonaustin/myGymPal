---
level: 128w
parent: auth.32w.md
deeper: auth-register.256w.md
relates-to:
  - ../controllers/auth.128w.md
source: views/auth/register.tpl
---

Renders the account creation page. Standalone HTML (own minimal navbar). A card holds username, password (min 8 chars),
and confirm-password fields. Client-side JS validates that the two password fields match using `setCustomValidity` and
adds Bootstrap `was-validated` class on submit.

Template variables: `.Error` (string, shown as a dismissible danger alert that auto-closes after 3 s), `.Username` (
pre-fills username on re-render after server-side error).

POSTs to `/register`. No AJAX. Registers `/sw.js` service worker.
