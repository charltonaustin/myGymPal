---
level: 256w
parent: root.32w.md
relates-to:
  - ../controllers/auth.128w.md
source: views/index.tpl
---

## Purpose

Public landing page that explains the application and guides new users through the setup flow.

## Partials Included

- `partials/navbar.tpl` (receives `.LoggedIn` and `.ActivePage` for conditional rendering)

## Template Variables

None beyond those consumed by the navbar partial. The page body contains no `{{.Var}}` substitutions.

## Content Sections

1. **Quick Start** — Two-step path: Create a Program, Start a workout. Includes a link to `/example`.
2. **Full Setup** — Four-step path: Create exercises (`/exercises/new`), Create a template (`/templates/new`), Create a
   program (`/programs/new`), Start a workout.
3. **Tracking Weight** — Explains daily weigh-in flow: `/weight` page, 3-day rolling average on dashboard, unit
   preference in `/settings`.
4. **Tracking Macros** — Explains food logging: partial macro tracking, automatic calorie calculation (4 kcal/g
   protein & carbs, 9 kcal/g fat), daily goals comparison, serving sizes.

## JavaScript Behavior

Loads `offline-sync.js` and registers `/sw.js` PWA service worker. No other JS.

## Navigation Links

`/example`, `/programs/new`, `/exercises/new`, `/templates/new`, `/weight`, `/settings`, `/macros`

## Flash Messages

None.

## AJAX / Fetch

None.
