package models

import (
	"errors"

	"github.com/beego/beego/v2/client/orm"
)

type Phase struct {
	ID          int64 `orm:"column(id);auto;pk"`
	ProgramID   int64 `orm:"column(program_id)"`
	PhaseNumber int   `orm:"column(phase_number)"`
	RepMin      int   `orm:"column(rep_min)"`
	RepMax      int   `orm:"column(rep_max)"`
	DefaultSets int   `orm:"column(default_sets)"`
	RestSeconds int   `orm:"column(rest_seconds)"`
}

func (p *Phase) TableName() string {
	return "phases"
}

func init() {
	orm.RegisterModel(&Phase{})
}

func GetPhasesByProgramID(programID int64) ([]*Phase, error) {
	o := orm.NewOrm()
	var phases []*Phase
	_, err := o.QueryTable(&Phase{}).Filter("ProgramID", programID).OrderBy("PhaseNumber").All(&phases)
	return phases, err
}

type PhaseUpdate struct {
	PhaseNumber int
	RepMin      int
	RepMax      int
	DefaultSets int
	RestSeconds int
}

func UpdatePhaseRepRanges(programID int64, updates []PhaseUpdate) error {
	for _, u := range updates {
		if u.RepMin <= 0 {
			return errors.New("rep_min must be greater than 0")
		}
		if u.RepMax < u.RepMin {
			return errors.New("rep_max must be at least rep_min")
		}
		if u.DefaultSets <= 0 {
			return errors.New("default_sets must be greater than 0")
		}
	}

	o := orm.NewOrm()
	for _, u := range updates {
		_, err := o.QueryTable(&Phase{}).
			Filter("ProgramID", programID).
			Filter("PhaseNumber", u.PhaseNumber).
			Update(orm.Params{
				"rep_min":      u.RepMin,
				"rep_max":      u.RepMax,
				"default_sets": u.DefaultSets,
				"rest_seconds": u.RestSeconds,
			})
		if err != nil {
			return err
		}
	}
	return nil
}
