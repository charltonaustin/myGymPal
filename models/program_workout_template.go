package models

import "github.com/beego/beego/v2/client/orm"

type ProgramWorkoutTemplate struct {
	ID            int64 `orm:"column(id);auto;pk"`
	ProgramID     int64 `orm:"column(program_id)"`
	WorkoutNumber int   `orm:"column(workout_number)"`
	TemplateID    int64 `orm:"column(template_id)"`
}

func (t *ProgramWorkoutTemplate) TableName() string {
	return "program_workout_templates"
}

func init() {
	orm.RegisterModel(&ProgramWorkoutTemplate{})
}

func GetWorkoutTemplatesByProgram(programID int64) ([]*ProgramWorkoutTemplate, error) {
	o := orm.NewOrm()
	var ts []*ProgramWorkoutTemplate
	_, err := o.QueryTable(&ProgramWorkoutTemplate{}).Filter("ProgramID", programID).All(&ts)
	return ts, err
}

func UpsertWorkoutTemplate(programID int64, workoutNumber int, templateID int64) error {
	o := orm.NewOrm()
	_, err := o.Raw(`
		INSERT INTO program_workout_templates (program_id, workout_number, template_id)
		VALUES (?, ?, ?)
		ON CONFLICT (program_id, workout_number)
		DO UPDATE SET template_id = EXCLUDED.template_id
	`, programID, workoutNumber, templateID).Exec()
	return err
}

func DeleteWorkoutTemplate(programID int64, workoutNumber int) error {
	o := orm.NewOrm()
	_, err := o.Raw(
		`DELETE FROM program_workout_templates WHERE program_id = ? AND workout_number = ?`,
		programID, workoutNumber,
	).Exec()
	return err
}
