package models

type ormSessionExerciseRepository struct{}

func NewSessionExerciseRepository() SessionExerciseRepository {
	return &ormSessionExerciseRepository{}
}

func (r *ormSessionExerciseRepository) Create(sessionID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string) (*SessionExercise, error) {
	return CreateSessionExercise(sessionID, name, isBodyweight, goalWeight, weightUnit)
}

func (r *ormSessionExerciseRepository) GetBySession(sessionID int64) ([]*SessionExerciseView, error) {
	return GetSessionExercisesWithSets(sessionID)
}

func (r *ormSessionExerciseRepository) GetByID(exerciseID int64) (*SessionExercise, error) {
	return GetSessionExerciseByID(exerciseID)
}

func (r *ormSessionExerciseRepository) LogSet(exerciseID int64, setNumber int, actualWeight float64, weightUnit string, actualReps int) (*SessionSet, error) {
	return LogSessionSet(exerciseID, setNumber, actualWeight, weightUnit, actualReps)
}

func (r *ormSessionExerciseRepository) CountSetsByExercise(exerciseID int64) (int, error) {
	return CountSetsByExercise(exerciseID)
}
