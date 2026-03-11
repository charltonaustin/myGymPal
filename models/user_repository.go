package models

type ormUserRepository struct{}

func NewUserRepository() UserRepository {
	return &ormUserRepository{}
}

func (r *ormUserRepository) Create(username, password, weightUnit string) (*User, error) {
	return CreateUser(username, password, weightUnit)
}

func (r *ormUserRepository) GetByUsername(username string) (*User, error) {
	return GetUserByUsername(username)
}

func (r *ormUserRepository) GetByID(id int64) (*User, error) {
	return GetUserByID(id)
}

func (r *ormUserRepository) UpdateWeightUnit(userID int64, unit string) error {
	return UpdateWeightUnit(userID, unit)
}

func (r *ormUserRepository) DeleteByUsername(username string) error {
	return DeleteUserByUsername(username)
}

func (r *ormUserRepository) DeleteByID(id int64) error {
	return DeleteUserByID(id)
}
