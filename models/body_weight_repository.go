package models

import "time"

type ormBodyWeightRepository struct{}

func NewBodyWeightRepository() BodyWeightRepository {
	return &ormBodyWeightRepository{}
}

func (r *ormBodyWeightRepository) Create(userID int64, date time.Time, weight float64, weightUnit string) (*BodyWeight, error) {
	return CreateBodyWeight(userID, date, weight, weightUnit)
}

func (r *ormBodyWeightRepository) GetAllByUser(userID int64) ([]*BodyWeight, error) {
	return GetBodyWeightsByUser(userID)
}

func (r *ormBodyWeightRepository) GetByID(id, userID int64) (*BodyWeight, error) {
	return GetBodyWeightByID(id, userID)
}

func (r *ormBodyWeightRepository) Update(id, userID int64, weight float64, weightUnit string) (*BodyWeight, error) {
	return UpdateBodyWeight(id, userID, weight, weightUnit)
}

func (r *ormBodyWeightRepository) Delete(id, userID int64) error {
	return DeleteBodyWeight(id, userID)
}
