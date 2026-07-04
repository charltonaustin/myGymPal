---
level: 256w
parent: partials.32w.md
relates-to:
  - ../controllers/auth.128w.md
source: views/partials/navbar.tpl
---

## Purpose

Shared navigation component included at the top of every authenticated page. Provides consistent site navigation and
login/logout controls.

## Template Variables

| Variable      | Type   | Description                                                                                                                           |
|---------------|--------|---------------------------------------------------------------------------------------------------------------------------------------|
| `.LoggedIn`   | bool   | Controls which nav links and action buttons are rendered                                                                              |
| `.ActivePage` | string | Slug of the current page (e.g. "programs", "exercises", "dashboard"); used to add `active fw-semibold` class to the matching nav link |

## Conditional Rendering

- `{{if .LoggedIn}}` — renders full nav (Home, Programs, Templates, Exercises, Weight, Macros, Dashboard, Settings + Log
  Out button); else renders Home link only + Log In and Create Account buttons.
- Each nav link: `{{if eq .ActivePage "X"}}active fw-semibold{{end}}` pattern applied to class attribute.
- `{{if isDev}}` — renders a red centered "DEVELOPMENT" banner `<div>` immediately below the `<nav>` element. `isDev` is
  a custom template function registered in `routers/router.go`.

## Structure

- `<nav class="navbar navbar-dark bg-dark navbar-expand-md">` — Bootstrap 5 dark navbar.
- `.navbar-toggler` — hamburger button for `xs`/`sm` screens, targets `#mainNav`.
- `#mainNav` — collapsible div; contains `div.navbar-nav.me-auto` (page links) and `div.d-flex` (action buttons).

## Navigation Links (authenticated)

`/` (Home), `/programs`, `/templates`, `/exercises`, `/weight`, `/macros`, `/dashboard`, `/settings`

## Action Buttons

- Authenticated: "Log Out" → `/logout` (`btn-outline-light`)
- Unauthenticated: "Log In" → `/login` (`btn-outline-light`), "Create Account" → `/register` (`btn-light`)

## JavaScript / AJAX

None. Pure HTML/Bootstrap.
