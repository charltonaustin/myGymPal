package models

import "time"

type ormSessionRepository struct{}

func NewSessionRepository() SessionRepository {
	return &ormSessionRepository{}
}

func (r *ormSessionRepository) Create(programID, userID int64, phaseNumber, weekNumber, workoutNumber int, isDeload bool, date time.Time) (*Session, error) {
	return CreateSession(programID, userID, phaseNumber, weekNumber, workoutNumber, isDeload, date)
}

func (r *ormSessionRepository) CountByProgram(programID int64) (int, error) {
	return CountSessionsByProgram(programID)
}

func (r *ormSessionRepository) LatestByProgram(programID int64) (*Session, error) {
	return GetLatestSessionByProgram(programID)
}

func (r *ormSessionRepository) GetByID(id, userID int64) (*Session, error) {
	return GetSessionByID(id, userID)
}

func (r *ormSessionRepository) GetByProgram(programID int64) ([]*Session, error) {
	return GetSessionsByProgram(programID)
}

func (r *ormSessionRepository) GetRecentByUser(userID int64, limit int) ([]*RecentSession, error) {
	return GetRecentSessionsByUser(userID, limit)
}

func (r *ormSessionRepository) Delete(id, userID int64) error {
	return DeleteSession(id, userID)
}
