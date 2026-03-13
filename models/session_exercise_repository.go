package models

type ormSessionExerciseRepository struct{}

func NewSessionExerciseRepository() SessionExerciseRepository {
	return &ormSessionExerciseRepository{}
}

func (r *ormSessionExerciseRepository) Create(sessionID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, goalReps int, block string) (*SessionExercise, error) {
	return CreateSessionExercise(sessionID, name, isBodyweight, goalWeight, weightUnit, goalReps, block)
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

func (r *ormSessionExerciseRepository) DeleteSet(setID int64) error {
	return DeleteSessionSet(setID)
}

func (r *ormSessionExerciseRepository) LogCardio(sessionExerciseID int64, cardioType string, goalDuration, actualDuration int) (*CardioLog, error) {
	return LogCardioEntry(sessionExerciseID, cardioType, goalDuration, actualDuration)
}

func (r *ormSessionExerciseRepository) DeleteCardioLog(id int64) error {
	return DeleteCardioLog(id)
}

func (r *ormSessionExerciseRepository) DeleteExercise(exerciseID int64) error {
	return DeleteSessionExercise(exerciseID)
}
