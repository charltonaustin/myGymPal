package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Program struct {
	ID              int64     `orm:"column(id);auto;pk"`
	UserID          int64     `orm:"column(user_id)"`
	Name            string    `orm:"column(name)"`
	StartDate       time.Time `orm:"column(start_date);type(date)"`
	NumPhases       int       `orm:"column(num_phases)"`
	WeeksPerPhase   int       `orm:"column(weeks_per_phase)"`
	WorkoutsPerWeek int       `orm:"column(workouts_per_week)"`
	CreatedAt       time.Time `orm:"column(created_at);auto_now_add"`
	UpdatedAt       time.Time `orm:"column(updated_at);auto_now"`
}

func (p *Program) TableName() string {
	return "programs"
}

func init() {
	orm.RegisterModel(&Program{})
}

func CreateProgram(userID int64, name string, startDate time.Time, numPhases, weeksPerPhase, workoutsPerWeek, defaultRepMin, defaultRepMax int) (*Program, error) {
	if name == "" {
		return nil, errors.New("program name is required")
	}
	if numPhases <= 0 {
		return nil, errors.New("num_phases must be greater than 0")
	}
	if weeksPerPhase <= 0 {
		return nil, errors.New("weeks_per_phase must be greater than 0")
	}
	if workoutsPerWeek <= 0 {
		return nil, errors.New("workouts_per_week must be greater than 0")
	}
	if defaultRepMin <= 0 {
		return nil, errors.New("default_rep_min must be greater than 0")
	}
	if defaultRepMax < defaultRepMin {
		return nil, errors.New("default_rep_max must be at least default_rep_min")
	}

	p := &Program{
		UserID:          userID,
		Name:            name,
		StartDate:       startDate,
		NumPhases:       numPhases,
		WeeksPerPhase:   weeksPerPhase,
		WorkoutsPerWeek: workoutsPerWeek,
	}

	tx, err := orm.NewOrm().Begin()
	if err != nil {
		return nil, err
	}

	if _, err := tx.Insert(p); err != nil {
		tx.Rollback()
		return nil, err
	}

	for i := 1; i <= numPhases; i++ {
		ph := &Phase{ProgramID: p.ID, PhaseNumber: i, RepMin: defaultRepMin, RepMax: defaultRepMax}
		if _, err := tx.Insert(ph); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return p, nil
}

func GetProgramsByUserID(userID int64) ([]*Program, error) {
	o := orm.NewOrm()
	var programs []*Program
	_, err := o.QueryTable(&Program{}).Filter("UserID", userID).OrderBy("StartDate").All(&programs)
	return programs, err
}

func GetProgramByID(id, userID int64) (*Program, error) {
	o := orm.NewOrm()
	p := &Program{ID: id}
	if err := o.Read(p); err != nil {
		return nil, err
	}
	if p.UserID != userID {
		return nil, errors.New("not found")
	}
	return p, nil
}

func DeleteProgram(id, userID int64) error {
	o := orm.NewOrm()
	n, err := o.QueryTable(&Program{}).Filter("ID", id).Filter("UserID", userID).Delete()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("not found")
	}
	return nil
}
