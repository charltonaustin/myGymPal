---
level: 128w
parent: ../controllers.32w.md
deeper: account.256w.md
relates-to:
  - auth.128w.md
source: controllers/account.go
---

The account controller manages user account settings. GET /settings requires a session; it fetches the user via
`Users.GetByID`, reads any pending flash success, and renders `settings.tpl` with the user's current `WeightUnit`. POST
/settings updates the weight unit via `Users.UpdateWeightUnit`; on success it stores a flash success message ("Settings
saved.") and redirects back to /settings. On failure it re-renders the form with an inline error. POST /account/delete
calls `Users.DeleteByID` to permanently remove the account; on success it destroys the session and redirects to /login.
On delete failure it re-renders `settings.tpl` with an error message and tries to reload the user's WeightUnit for
display. All three handlers redirect unauthenticated requests to /login.
