package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Template struct {
	ID        int64     `orm:"column(id);auto;pk"`
	Name      string    `orm:"column(name)"`
	Focus     string    `orm:"column(focus)"`
	CreatedAt time.Time `orm:"column(created_at);auto_now_add"`
	UpdatedAt time.Time `orm:"column(updated_at);auto_now"`
}

func (t *Template) TableName() string {
	return "templates"
}

type TemplateExercise struct {
	ID           int64   `orm:"column(id);auto;pk"`
	TemplateID   int64   `orm:"column(template_id)"`
	Name         string  `orm:"column(name)"`
	IsBodyweight bool    `orm:"column(is_bodyweight)"`
	GoalWeight   float64 `orm:"column(goal_weight)"`
	WeightUnit   string  `orm:"column(weight_unit)"`
	RepMin       int     `orm:"column(rep_min)"`
	RepMax       int     `orm:"column(rep_max)"`
	SortOrder    int     `orm:"column(sort_order)"`
}

func (e *TemplateExercise) TableName() string {
	return "template_exercises"
}

type TemplateExerciseInput struct {
	Name         string
	IsBodyweight bool
	GoalWeight   float64
	WeightUnit   string
	RepMin       int
	RepMax       int
	SortOrder    int
}

func init() {
	orm.RegisterModel(&Template{}, &TemplateExercise{})
}

func CreateTemplate(name, focus string, exercises []TemplateExerciseInput) (*Template, error) {
	if name == "" {
		return nil, errors.New("template name is required")
	}
	if len(exercises) == 0 {
		return nil, errors.New("at least one exercise is required")
	}
	for _, ex := range exercises {
		if ex.Name == "" {
			return nil, errors.New("exercise name is required")
		}
		if ex.RepMin <= 0 {
			return nil, errors.New("rep_min must be greater than 0")
		}
		if ex.RepMax < ex.RepMin {
			return nil, errors.New("rep_max must be at least rep_min")
		}
		if !ex.IsBodyweight && ex.GoalWeight < 0 {
			return nil, errors.New("goal weight must be 0 or greater")
		}
	}

	t := &Template{Name: name, Focus: focus}

	tx, err := orm.NewOrm().Begin()
	if err != nil {
		return nil, err
	}

	if _, err := tx.Insert(t); err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, ex := range exercises {
		e := &TemplateExercise{
			TemplateID:   t.ID,
			Name:         ex.Name,
			IsBodyweight: ex.IsBodyweight,
			GoalWeight:   ex.GoalWeight,
			WeightUnit:   ex.WeightUnit,
			RepMin:       ex.RepMin,
			RepMax:       ex.RepMax,
			SortOrder:    ex.SortOrder,
		}
		if _, err := tx.Insert(e); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return t, nil
}

func GetAllTemplates() ([]*Template, error) {
	o := orm.NewOrm()
	var templates []*Template
	_, err := o.QueryTable(&Template{}).OrderBy("Name").All(&templates)
	return templates, err
}

func GetTemplateByID(id int64) (*Template, []*TemplateExercise, error) {
	o := orm.NewOrm()
	t := &Template{ID: id}
	if err := o.Read(t); err != nil {
		return nil, nil, errors.New("not found")
	}
	var exercises []*TemplateExercise
	_, err := o.QueryTable(&TemplateExercise{}).Filter("TemplateID", id).OrderBy("SortOrder").All(&exercises)
	if err != nil {
		return nil, nil, err
	}
	return t, exercises, nil
}
