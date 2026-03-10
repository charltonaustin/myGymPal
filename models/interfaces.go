package models

import "time"

type UserRepository interface {
	Create(username, password, weightUnit string) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByID(id int64) (*User, error)
	UpdateWeightUnit(userID int64, unit string) error
	DeleteByUsername(username string) error
}

type ProgramRepository interface {
	Create(userID int64, name string, startDate time.Time, numPhases, weeksPerPhase, defaultRepMin, defaultRepMax int) (*Program, error)
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
	GetAll() ([]*Template, error)
	GetByID(id int64) (*Template, []*TemplateExercise, error)
}
