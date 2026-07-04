---
level: 128w
parent: ../models.32w.md
deeper: user.256w.md
relates-to:
  - ../controllers/auth.128w.md
  - ../controllers/account.128w.md
source: models/user.go, models/user_repository.go
---

# User

`User` represents an authenticated application user. It stores the login `Username` (unique), a bcrypt `PasswordHash`,
and a `WeightUnit` preference (`"lb"` or `"kg"`). Auto-managed `CreatedAt` and `UpdatedAt` timestamps are set by the
ORM.

Password verification is a method on the struct (`CheckPassword`), keeping bcrypt logic contained. `UpdateWeightUnit`
validates that the unit is strictly `"lb"` or `"kg"` before persisting.

The `UserRepository` interface exposes six methods: `Create`, `GetByUsername`, `GetByID`, `UpdateWeightUnit`,
`DeleteByUsername`, and `DeleteByID`. The ORM implementation in `user_repository.go` delegates to package-level
functions in `user.go`.

This model is central to session-based authentication: the auth controller stores `user_id` in the session after
verifying credentials, and every other repository query is scoped to `userID`. The account controller uses
`UpdateWeightUnit` to handle the preferences form.
