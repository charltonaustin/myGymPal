---
level: 256w
parent: auth.32w.md
relates-to:
  - ../controllers/auth.128w.md
source: views/auth/login.tpl
---

## Purpose

Standalone login page. Does not include the shared navbar partial — it renders its own minimal dark navbar with only the
brand link.

## Template Variables

| Variable    | Type   | Description                                                                                                 |
|-------------|--------|-------------------------------------------------------------------------------------------------------------|
| `.Error`    | string | Error message displayed in a `div.alert.alert-danger` when login fails. Conditional on `{{if .Error}}`.     |
| `.Username` | string | Pre-fills the username input `value` attribute so the user does not need to retype it after a failed login. |

## Form Fields

Form: `method="POST" action="/login" novalidate`

| Field    | Name       | Type     | Notes                                                              |
|----------|------------|----------|--------------------------------------------------------------------|
| Username | `username` | text     | `autocomplete="username"`, `required`, pre-filled with `.Username` |
| Password | `password` | password | `autocomplete="current-password"`, `required`                      |
| Submit   | —          | submit   | "Log In", full-width dark button                                   |

## Conditional Rendering

- `.Error` block: renders a `div.alert.alert-danger` containing the error string.

## JavaScript Behavior

No form-validation JavaScript. The page loads Bootstrap bundle (CDN), `/static/offline-sync.js`, and registers `/sw.js`
as the PWA service worker.

## Flash Messages

None. Error feedback is passed directly as `.Error` from the controller, not via Beego flash.

## Navigation

Link to `/register` below the form for users without accounts.
