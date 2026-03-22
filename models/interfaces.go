package models

import "time"

type UserRepository interface {
	Create(username, password, weightUnit string) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByID(id int64) (*User, error)
	UpdateWeightUnit(userID int64, unit string) error
	DeleteByUsername(username string) error
	DeleteByID(id int64) error
}

type ProgramRepository interface {
	Create(userID int64, name string, startDate time.Time, numPhases, weeksPerPhase, workoutsPerWeek, defaultRepMin, defaultRepMax, defaultSets int) (*Program, error)
	GetAllByUser(userID int64) ([]*Program, error)
	GetByID(id, userID int64) (*Program, error)
	Delete(id, userID int64) error
}

type PhaseRepository interface {
	GetByProgram(programID int64) ([]*Phase, error)
	UpdateRepRanges(programID int64, updates []PhaseUpdate) error
}

type TemplateRepository interface {
	Create(name, focus string, exercises []TemplateExerciseInput) (*Template, error)
	Update(id int64, name, focus string, exercises []TemplateExerciseInput) (*Template, error)
	GetAll() ([]*Template, error)
	GetByID(id int64) (*Template, []*TemplateExercise, error)
	Delete(id int64) error
}

type SessionRepository interface {
	Create(programID, userID int64, phaseNumber, weekNumber, workoutNumber int, isDeload bool, date time.Time) (*Session, error)
	CountByProgram(programID int64) (int, error)
	LatestByProgram(programID int64) (*Session, error)
	GetByID(id, userID int64) (*Session, error)
	GetByProgram(programID int64) ([]*Session, error)
	GetRecentByUser(userID int64, limit int) ([]*RecentSession, error)
	Delete(id, userID int64) error
}

type ExerciseRepository interface {
	Create(userID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, isTimeBased bool, goalSeconds int, goalRepMin int, goalRepMax int) (*Exercise, error)
	GetAllByUser(userID int64) ([]*Exercise, error)
	GetByID(id, userID int64) (*Exercise, error)
	GetByName(userID int64, name string) (*Exercise, error)
	Update(id, userID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, isTimeBased bool, goalSeconds int, goalRepMin int, goalRepMax int) (*Exercise, error)
	UpdateGoalWeight(id int64, goalWeight float64) error
	Delete(id, userID int64) error
}

type BodyWeightRepository interface {
	Create(userID int64, date time.Time, weight float64, weightUnit string) (*BodyWeight, error)
	GetAllByUser(userID int64) ([]*BodyWeight, error)
	GetByID(id, userID int64) (*BodyWeight, error)
	Update(id, userID int64, weight float64, weightUnit string) (*BodyWeight, error)
	Delete(id, userID int64) error
}

type MacroGoalRepository interface {
	Get(userID int64) (*MacroGoal, error)
	Upsert(userID int64, protein, carbs, fat float64) (*MacroGoal, error)
}

type SessionExerciseRepository interface {
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
}
