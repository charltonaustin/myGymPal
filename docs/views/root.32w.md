---
level: 32w
parent: ../views.32w.md
deeper: root-index.128w.md
relates-to:
  - ../controllers/auth.128w.md
source: views/
---

| Template  | Summary                                                                                                                                                     | Detail                         |
|-----------|-------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------|
| index     | Marketing/landing page with Quick Start, Full Setup, Weight Tracking, and Macros Tracking sections; links throughout; no auth required                      | [128w](root-index.128w.md)     |
| dashboard | Authenticated home: weight N-day average card, macros N-day average card (with goal % when set), recent workout list with link to next session              | [128w](root-dashboard.128w.md) |
| settings  | Account Settings page: lb/kg weight-unit preference (POSTs to `/settings`); Danger Zone with delete-account confirmation modal (POSTs to `/account/delete`) | [128w](root-settings.128w.md)  |
| error     | Generic error fallback page: static "Something went wrong" message with link home; no template variables                                                    | [128w](root-error.128w.md)     |
| offline   | PWA offline fallback: "You're Offline" message with Retry button (`window.location.reload()`); shown by service worker when no cached response              | [128w](root-offline.128w.md)   |
