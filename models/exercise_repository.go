package models

type ormExerciseRepository struct{}

func NewExerciseRepository() ExerciseRepository {
	return &ormExerciseRepository{}
}

func (r *ormExerciseRepository) Create(userID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, isTimeBased bool, goalSeconds int, goalRepMin int, goalRepMax int) (*Exercise, error) {
	return CreateExercise(userID, name, isBodyweight, goalWeight, weightUnit, isTimeBased, goalSeconds, goalRepMin, goalRepMax)
}

func (r *ormExerciseRepository) GetAllByUser(userID int64) ([]*Exercise, error) {
	return GetExercisesByUser(userID)
}

func (r *ormExerciseRepository) GetByID(id, userID int64) (*Exercise, error) {
	return GetExerciseByID(id, userID)
}

func (r *ormExerciseRepository) GetByName(userID int64, name string) (*Exercise, error) {
	return GetExerciseByName(userID, name)
}

func (r *ormExerciseRepository) Update(id, userID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, isTimeBased bool, goalSeconds int, goalRepMin int, goalRepMax int) (*Exercise, error) {
	return UpdateExercise(id, userID, name, isBodyweight, goalWeight, weightUnit, isTimeBased, goalSeconds, goalRepMin, goalRepMax)
}

func (r *ormExerciseRepository) UpdateGoalWeight(id int64, goalWeight float64) error {
	return UpdateExerciseGoalWeight(id, goalWeight)
}

func (r *ormExerciseRepository) Delete(id, userID int64) error {
	return DeleteExercise(id, userID)
}
