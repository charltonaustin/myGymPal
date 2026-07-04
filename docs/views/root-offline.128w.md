---
level: 128w
parent: root.32w.md
deeper: root-offline.256w.md
relates-to: []
source: views/offline.tpl
---

PWA offline fallback page served by the service worker when a navigation request fails and no cached response exists.
Standalone HTML with a minimal dark navbar (brand only). Shows a signal-bars emoji, "You're Offline" heading, and a "
Retry" button that calls `window.location.reload()`. No template variables. No AJAX. No service-worker registration (
this page is itself the fallback).
