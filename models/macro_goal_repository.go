package models

type ormMacroGoalRepository struct{}

func NewMacroGoalRepository() MacroGoalRepository {
	return &ormMacroGoalRepository{}
}

func (r *ormMacroGoalRepository) Get(userID int64) (*MacroGoal, error) {
	return GetMacroGoal(userID)
}

func (r *ormMacroGoalRepository) Upsert(userID int64, protein, carbs, fat float64) (*MacroGoal, error) {
	return UpsertMacroGoal(userID, protein, carbs, fat)
}
