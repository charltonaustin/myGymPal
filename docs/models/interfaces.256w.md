---
level: 256w
parent: ../models.32w.md
relates-to:
  - user.128w.md
  - program.128w.md
  - phase.128w.md
  - session.128w.md
  - session_exercise.128w.md
  - template.128w.md
  - exercise.128w.md
  - macro_goal.128w.md
  - body_weight.128w.md
source: models/interfaces.go
---

# Interfaces (full reference)

`models/interfaces.go` declares the Go interfaces that define the repository contract for each domain model. All except
`MacroRepository` are here.

## UserRepository

```go
Create(username, password, weightUnit string) (*User, error)
GetByUsername(username string) (*User, error)
GetByID(id int64) (*User, error)
UpdateWeightUnit(userID int64, unit string) error
DeleteByUsername(username string) error
DeleteByID(id int64) error
```

## ProgramRepository

```go
Create(userID int64, name string, startDate time.Time, numPhases, weeksPerPhase, workoutsPerWeek, defaultRepMin, defaultRepMax, defaultSets int) (*Program, error)
GetAllByUser(userID int64) ([]*Program, error)
GetByID(id, userID int64) (*Program, error)
Delete(id, userID int64) error
```

## PhaseRepository

```go
GetByProgram(programID int64) ([]*Phase, error)
UpdateRepRanges(programID int64, updates []PhaseUpdate) error
```

## TemplateRepository

```go
Create(name, focus string, circuits []TemplateCircuitInput, exercises []TemplateExerciseInput) (*Template, error)
Update(id int64, name, focus string, circuits []TemplateCircuitInput, exercises []TemplateExerciseInput) (*Template, error)
GetAll() ([]*Template, error)
GetByID(id int64) (*Template, []*TemplateExercise, error)
GetCircuits(templateID int64) ([]*TemplateCircuit, error)
Delete(id int64) error
```

`GetByID` returns exercises but not circuits: the session and exercise controllers call it for the exercise list
alone. Circuits are a separate read so those callers stay unaffected.

## SessionRepository

```go
Create(programID, userID int64, phaseNumber, weekNumber, workoutNumber int, isDeload bool, date time.Time) (*Session, error)
CountByProgram(programID int64) (int, error)
LatestByProgram(programID int64) (*Session, error)
GetByID(id, userID int64) (*Session, error)
GetByProgram(programID int64) ([]*Session, error)
GetRecentByUser(userID int64, limit int) ([]*RecentSession, error)
Delete(id, userID int64) error
```

## ExerciseRepository

```go
Create(userID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, isTimeBased bool, goalSeconds int, goalRepMin int, goalRepMax int, defaultBlock string) (*Exercise, error)
GetAllByUser(userID int64) ([]*Exercise, error)
GetByID(id, userID int64) (*Exercise, error)
GetByName(userID int64, name string) (*Exercise, error)
Update(id, userID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, isTimeBased bool, goalSeconds int, goalRepMin int, goalRepMax int, defaultBlock string) (*Exercise, error)
UpdateGoalWeight(id int64, goalWeight float64, weightUnit string) error
Delete(id, userID int64) error
```

## BodyWeightRepository

```go
Create(userID int64, date time.Time, weight float64, weightUnit string) (*BodyWeight, error)
GetAllByUser(userID int64) ([]*BodyWeight, error)
GetByID(id, userID int64) (*BodyWeight, error)
Update(id, userID int64, weight float64, weightUnit string) (*BodyWeight, error)
Delete(id, userID int64) error
```

## MacroGoalRepository

```go
Get(userID int64) (*MacroGoal, error)
Upsert(userID int64, protein, carbs, fat float64) (*MacroGoal, error)
```

## SessionExerciseRepository

```go
Create(sessionID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, goalReps int, block string, isTimeBased bool, goalSeconds int) (*SessionExercise, error)
GetBySession(sessionID int64) ([]*SessionExerciseView, error)
GetByID(exerciseID int64) (*SessionExercise, error)
LogSet(exerciseID int64, setNumber int, actualWeight float64, weightUnit string, actualReps int, actualSeconds int, activityType string) (*SessionSet, error)
CountSetsByExercise(exerciseID int64) (int, error)
DeleteSet(setID int64) error
LogCardio(sessionExerciseID int64, cardioType string, goalDuration, actualDuration int) (*CardioLog, error)
DeleteCardioLog(id int64) error
DeleteExercise(exerciseID int64) error
UpdateSortOrders(sessionID int64, ids []int64) error
```

## Not in this file

`MacroRepository` (for `MacroEntry`) is defined in `models/macro_entry_repository.go`.

## Usage pattern

- Package-level variables in `controllers/repos.go` hold interface values.
- ORM implementations injected in `routers/router.go:Register()`.
- Test mocks in `controllers/mocks_test.go` implement these interfaces with function fields.
