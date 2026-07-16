package models

type ormExerciseRepository struct{}

func NewExerciseRepository() ExerciseRepository {
	return &ormExerciseRepository{}
}

func (r *ormExerciseRepository) Create(userID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, isTimeBased bool, goalSeconds int, goalRepMin int, goalRepMax int, defaultBlock string) (*Exercise, error) {
	return CreateExercise(userID, name, isBodyweight, goalWeight, weightUnit, isTimeBased, goalSeconds, goalRepMin, goalRepMax, defaultBlock)
}

func (r *ormExerciseRepository) GetAll(userID int64) ([]*Exercise, error) {
	return GetExercisesAll(userID)
}

func (r *ormExerciseRepository) GetAllByUser(userID int64) ([]*Exercise, error) {
	return GetExercisesByUser(userID)
}

func (r *ormExerciseRepository) GetAvailableGlobalNames(userID int64) ([]string, error) {
	return GetGlobalExercisesNotConfigured(userID)
}

func (r *ormExerciseRepository) GetByID(id, userID int64) (*Exercise, error) {
	return GetExerciseByID(id, userID)
}

func (r *ormExerciseRepository) GetByName(userID int64, name string) (*Exercise, error) {
	return GetExerciseByName(userID, name)
}

func (r *ormExerciseRepository) Update(id, userID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, isTimeBased bool, goalSeconds int, goalRepMin int, goalRepMax int, defaultBlock string) (*Exercise, error) {
	return UpdateExercise(id, userID, name, isBodyweight, goalWeight, weightUnit, isTimeBased, goalSeconds, goalRepMin, goalRepMax, defaultBlock)
}

func (r *ormExerciseRepository) UpdateGoalWeight(id, userID int64, goalWeight float64, weightUnit string) error {
	return UpdateExerciseGoalWeight(id, userID, goalWeight, weightUnit)
}

func (r *ormExerciseRepository) Delete(id, userID int64) error {
	return DeleteExercise(id, userID)
}

func (r *ormExerciseRepository) GetHistory(userID int64, names []string, unit string, days int) ([]ExerciseHistorySeries, error) {
	return GetExerciseHistory(userID, names, unit, days)
}

func (r *ormExerciseRepository) GetRecentNames(userID int64, days int) ([]string, error) {
	return GetRecentExerciseNames(userID, days)
}
