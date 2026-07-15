package models

type ormProgramWorkoutTemplateRepository struct{}

func NewProgramWorkoutTemplateRepository() ProgramWorkoutTemplateRepository {
	return &ormProgramWorkoutTemplateRepository{}
}

func (r *ormProgramWorkoutTemplateRepository) GetByProgram(programID int64) ([]*ProgramWorkoutTemplate, error) {
	return GetWorkoutTemplatesByProgram(programID)
}

func (r *ormProgramWorkoutTemplateRepository) Upsert(programID int64, workoutNumber int, templateID int64) error {
	return UpsertWorkoutTemplate(programID, workoutNumber, templateID)
}

func (r *ormProgramWorkoutTemplateRepository) Delete(programID int64, workoutNumber int) error {
	return DeleteWorkoutTemplate(programID, workoutNumber)
}
