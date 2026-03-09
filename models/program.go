package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Program struct {
	ID        int64     `orm:"column(id);auto;pk"`
	UserID    int64     `orm:"column(user_id)"`
	Name      string    `orm:"column(name)"`
	StartDate time.Time `orm:"column(start_date);type(date)"`
	NumPhases int       `orm:"column(num_phases)"`
	CreatedAt time.Time `orm:"column(created_at);auto_now_add"`
	UpdatedAt time.Time `orm:"column(updated_at);auto_now"`
}

func (p *Program) TableName() string {
	return "programs"
}

func init() {
	orm.RegisterModel(&Program{})
}

func CreateProgram(userID int64, name string, startDate time.Time, numPhases int) (*Program, error) {
	if name == "" {
		return nil, errors.New("program name is required")
	}
	if numPhases <= 0 {
		return nil, errors.New("num_phases must be greater than 0")
	}

	p := &Program{
		UserID:    userID,
		Name:      name,
		StartDate: startDate,
		NumPhases: numPhases,
	}

	o := orm.NewOrm()
	if _, err := o.Insert(p); err != nil {
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
