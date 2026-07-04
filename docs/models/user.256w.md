---
level: 256w
parent: ../models.32w.md
relates-to:
  - ../controllers/auth.128w.md
  - ../controllers/account.128w.md
source: models/user.go, models/user_repository.go
---

# User (full reference)

## Struct fields

| Field        | Go type   | ORM tag / notes              |
|--------------|-----------|------------------------------|
| ID           | int64     | `auto;pk`                    |
| Username     | string    | `unique`                     |
| PasswordHash | string    | bcrypt hash stored as string |
| WeightUnit   | string    | `"lb"` or `"kg"`             |
| CreatedAt    | time.Time | `auto_now_add`               |
| UpdatedAt    | time.Time | `auto_now`                   |

## Repository interface (UserRepository)

```go
Create(username, password, weightUnit string) (*User, error)
GetByUsername(username string) (*User, error)
GetByID(id int64) (*User, error)
UpdateWeightUnit(userID int64, unit string) error
DeleteByUsername(username string) error
DeleteByID(id int64) error
```

## Notable behavior

- `Create` calls `bcrypt.GenerateFromPassword` at `bcrypt.DefaultCost` before inserting; returns an error if username or
  password is empty.
- `CheckPassword(password string) bool` is a struct method that wraps `bcrypt.CompareHashAndPassword`.
- `UpdateWeightUnit` validates the unit is `"lb"` or `"kg"` and returns an error otherwise; uses an ORM
  `QueryTable.Update` (not a full struct update).
- `GetByUsername` reads by the `Username` field; `GetByID` reads by primary key.
- Delete methods use `QueryTable.Filter.Delete` so no read is required first.

## ORM / SQL patterns

- `orm.RegisterModel(&User{})` called in `init()`.
- All mutations use the default ORM instance (`orm.NewOrm()`).
- No raw SQL; all queries go through Beego ORM's query builder.

## Relationships

- One user owns many Programs, Sessions, Exercises, MacroEntries, a MacroGoal, and BodyWeight entries — all scoped by
  `user_id` FK in those tables.

## Usage

Auth controller: login/register/logout. Account controller: weight-unit preference update. Every other controller reads
`user_id` from the session and passes it to scoped repositories.
