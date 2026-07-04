---
level: 256w
parent: auth.32w.md
relates-to:
  - ../controllers/auth.128w.md
source: views/auth/register.tpl
---

## Purpose

Account registration page. Standalone HTML with a minimal dark navbar (brand only). Collects username and password from
a new user.

## Template Variables

| Variable    | Type   | Description                                                                                                                                                                |
|-------------|--------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `.Error`    | string | Server-side error (e.g. username taken). Shown as `div.alert.alert-danger.alert-dismissible` with `id="error-alert"`. Auto-dismissed after 3 000 ms via `bootstrap.Alert`. |
| `.Username` | string | Pre-fills username input after server re-render on error so user does not retype it.                                                                                       |

## Form Fields

Form: `id="register-form" method="POST" action="/register" novalidate`

| Field            | Name               | Type     | Notes                                                      |
|------------------|--------------------|----------|------------------------------------------------------------|
| Username         | `username`         | text     | `autocomplete="username"`, `required`                      |
| Password         | `password`         | password | `autocomplete="new-password"`, `minlength="8"`, `required` |
| Confirm Password | `confirm_password` | password | `autocomplete="new-password"`, `required`                  |
| Submit           | —                  | submit   | "Create Account"                                           |

## JavaScript Behavior

- `checkPasswordMatch()`: calls `confirm.setCustomValidity(...)` — sets a validation message when the confirm field is
  non-empty and does not equal the password field; clears it otherwise.
- `password` and `confirm_password` both fire `checkPasswordMatch` on `input` events.
- Form `submit` handler: runs `checkPasswordMatch()`, prevents submission if `checkValidity()` is false, adds
  `was-validated` class for Bootstrap inline error display.
- Error alert auto-dismiss: `setTimeout(() => bootstrap.Alert.getInstance(alertEl).close(), 3000)`.

## Conditional Rendering

`.Error` block renders only when the string is non-empty.

## Flash Messages

None — error passed as `.Error` from controller.
