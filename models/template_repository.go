package models

type ormTemplateRepository struct{}

func NewTemplateRepository() TemplateRepository {
	return &ormTemplateRepository{}
}

func (r *ormTemplateRepository) Create(name, focus string, exercises []TemplateExerciseInput) (*Template, error) {
	return CreateTemplate(name, focus, exercises)
}

func (r *ormTemplateRepository) GetAll() ([]*Template, error) {
	return GetAllTemplates()
}

func (r *ormTemplateRepository) GetByID(id int64) (*Template, []*TemplateExercise, error) {
	return GetTemplateByID(id)
}

func (r *ormTemplateRepository) Update(id int64, name, focus string, exercises []TemplateExerciseInput) (*Template, error) {
	return UpdateTemplate(id, name, focus, exercises)
}

func (r *ormTemplateRepository) Delete(id int64) error {
	return DeleteTemplate(id)
}
