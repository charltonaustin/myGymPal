---
level: 128w
parent: ../models.32w.md
deeper: cardio_log.256w.md
relates-to:
  - session_exercise.128w.md
  - ../controllers/session.128w.md
source: models/cardio_log.go
---

# CardioLog

`CardioLog` is a cardio activity sub-record linked to a `SessionExercise` (which must have `Block == "cardio"`). It
stores `CardioType` (a string label), `GoalDuration` and `ActualDuration` (both in seconds as int), and a `CreatedAt`
auto-timestamp.

There is no standalone `CardioLogRepository` interface in `interfaces.go`. Cardio operations are exposed through
`SessionExerciseRepository`: `LogCardio(sessionExerciseID, cardioType, goalDuration, actualDuration)` and
`DeleteCardioLog(id)`.

`GetCardioLogsByExercise` (package-level function) fetches logs ordered by `CreatedAt`, but this function is not exposed
through any interface — it is called internally when `GetSessionExercisesWithSets` builds `SessionExerciseView` for
cardio-block exercises. The session controller is the sole consumer.
