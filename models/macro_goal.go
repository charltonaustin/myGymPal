package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type MacroGoal struct {
	ID        int64     `orm:"column(id);auto;pk"`
	UserID    int64     `orm:"column(user_id);unique"`
	Protein   float64   `orm:"column(protein)"`
	Carbs     float64   `orm:"column(carbs)"`
	Fat       float64   `orm:"column(fat)"`
	UpdatedAt time.Time `orm:"column(updated_at);auto_now"`
}

func (m *MacroGoal) TableName() string {
	return "macro_goals"
}

func init() {
	orm.RegisterModel(&MacroGoal{})
}

func GetMacroGoal(userID int64) (*MacroGoal, error) {
	o := orm.NewOrm()
	g := &MacroGoal{UserID: userID}
	err := o.Read(g, "UserID")
	if err == orm.ErrNoRows {
		return nil, nil
	}
	return g, err
}

func UpsertMacroGoal(userID int64, protein, carbs, fat float64) (*MacroGoal, error) {
	o := orm.NewOrm()
	g := &MacroGoal{UserID: userID}
	err := o.Read(g, "UserID")
	if err == orm.ErrNoRows {
		g.Protein = protein
		g.Carbs = carbs
		g.Fat = fat
		_, err = o.Insert(g)
		return g, err
	}
	if err != nil {
		return nil, err
	}
	g.Protein = protein
	g.Carbs = carbs
	g.Fat = fat
	_, err = o.Update(g, "protein", "carbs", "fat", "updated_at")
	return g, err
}
