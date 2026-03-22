package models

import "time"

type ormProgramRepository struct{}

func NewProgramRepository() ProgramRepository {
	return &ormProgramRepository{}
}

func (r *ormProgramRepository) Create(userID int64, name string, startDate time.Time, numPhases, weeksPerPhase, workoutsPerWeek, defaultRepMin, defaultRepMax, defaultSets int) (*Program, error) {
	return CreateProgram(userID, name, startDate, numPhases, weeksPerPhase, workoutsPerWeek, defaultRepMin, defaultRepMax, defaultSets)
}

func (r *ormProgramRepository) GetAllByUser(userID int64) ([]*Program, error) {
	return GetProgramsByUserID(userID)
}

func (r *ormProgramRepository) GetByID(id, userID int64) (*Program, error) {
	return GetProgramByID(id, userID)
}

func (r *ormProgramRepository) Delete(id, userID int64) error {
	return DeleteProgram(id, userID)
}
