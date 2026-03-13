package models

import "time"

type MacroRepository interface {
	Create(userID int64, date time.Time, foodName string, servingWeight float64, servingUnit string, protein, carbs, fat float64) (*MacroEntry, error)
	GetAllByUser(userID int64) ([]*MacroEntry, error)
	GetByID(id, userID int64) (*MacroEntry, error)
	Update(id, userID int64, foodName string, servingWeight float64, servingUnit string, protein, carbs, fat float64) (*MacroEntry, error)
	Delete(id, userID int64) error
}

type ormMacroRepository struct{}

func NewMacroRepository() MacroRepository {
	return &ormMacroRepository{}
}

func (r *ormMacroRepository) Create(userID int64, date time.Time, foodName string, servingWeight float64, servingUnit string, protein, carbs, fat float64) (*MacroEntry, error) {
	return CreateMacroEntry(userID, date, foodName, servingWeight, servingUnit, protein, carbs, fat)
}

func (r *ormMacroRepository) GetAllByUser(userID int64) ([]*MacroEntry, error) {
	return GetMacroEntriesByUser(userID)
}

func (r *ormMacroRepository) GetByID(id, userID int64) (*MacroEntry, error) {
	return GetMacroEntryByID(id, userID)
}

func (r *ormMacroRepository) Update(id, userID int64, foodName string, servingWeight float64, servingUnit string, protein, carbs, fat float64) (*MacroEntry, error) {
	return UpdateMacroEntry(id, userID, foodName, servingWeight, servingUnit, protein, carbs, fat)
}

func (r *ormMacroRepository) Delete(id, userID int64) error {
	return DeleteMacroEntry(id, userID)
}
