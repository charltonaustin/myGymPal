package models

type ormPhaseRepository struct{}

func NewPhaseRepository() PhaseRepository {
	return &ormPhaseRepository{}
}

func (r *ormPhaseRepository) GetByProgram(programID int64) ([]*Phase, error) {
	return GetPhasesByProgramID(programID)
}

func (r *ormPhaseRepository) UpdateRepRanges(programID int64, updates []PhaseUpdate) error {
	return UpdatePhaseRepRanges(programID, updates)
}
