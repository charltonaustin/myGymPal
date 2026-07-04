---
level: 256w
parent: ../controllers.32w.md
relates-to:
  - dashboard.256w.md
  - account.256w.md
source: controllers/weight.go
---

## Routes

| Method | Path               | Handler |
|--------|--------------------|---------|
| GET    | /weight            | Index   |
| POST   | /weight            | Create  |
| POST   | /weight/:id        | Update  |
| POST   | /weight/:id/delete | Delete  |

## Auth requirement

All handlers check `c.GetSession("user_id")`; nil redirects to /login.

## Session keys

- Read: `user_id` (int64)

## Template variables — Index

- `c.Data["LoggedIn"]` = true
- `c.Data["ActivePage"]` = "weight"
- `c.Data["Entries"]` = `[]*models.BodyWeight` (weights converted in-place to user's preferred unit)
- `c.Data["WeightUnit"]` = user's preferred unit string (from `Users.GetByID`)
- `c.Data["DefaultDate"]` = today in "2006-01-02" format
- `c.Data["WeightAvg"]` = `*weightAverage` (nil if no entries)
- `c.Data["Success"]` = flash success string (if present)

## Template

- `weight/index.tpl`

## Repository calls

- `BodyWeights.GetAllByUser(userID)` — Index
- `Users.GetByID(userID)` — Index (to read preferred unit)
- `BodyWeights.Create(userID, date, weight, weightUnit)` — Create
- `BodyWeights.Update(id, userID, weight, weightUnit)` — Update
- `BodyWeights.Delete(id, userID)` — Delete

## Helper function

`computeWeightAverage(entries, targetUnit)` — defined in `weight.go`, also called by `DashboardController.Get`. Averages
the most recent 3 entries (after unit conversion) into a `weightAverage{Days int, Weight float64, Unit string}` struct.
Returns nil if entries is empty.

## Unit handling

- Unit validation: if submitted `weight_unit` is not "kg", it is forced to "lb".
- Values stored in the original submitted unit.
- On Index, each entry is converted in-place: `e.Weight = models.ConvertWeight(e.Weight, e.WeightUnit, weightUnit)` then
  `e.WeightUnit = weightUnit`.
- Average is computed before in-place conversion so original units are passed to `computeWeightAverage`.

## Flash messages

Flash is read on Index (`flash.Data["success"]`), but no POST handler writes a flash — all POSTs simply redirect to
/weight on both success and failure (errors are logged).

## Redirect paths

- Create → /weight (always, errors only logged)
- Update → /weight (always, errors only logged)
- Delete → /weight (always, errors only logged)

## Relationship to other controllers

`computeWeightAverage` is re-used directly by `DashboardController.Get`. The preferred weight unit displayed here is
controlled by `AccountController.SettingsPost`.
