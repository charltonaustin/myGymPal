package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type MacroEntry struct {
	ID            int64     `orm:"column(id);auto;pk"`
	UserID        int64     `orm:"column(user_id)"`
	Date          time.Time `orm:"column(date);type(date)"`
	FoodName      string    `orm:"column(food_name)"`
	ServingWeight float64   `orm:"column(serving_weight)"`
	ServingUnit   string    `orm:"column(serving_unit)"`
	Protein       float64   `orm:"column(protein)"`
	Carbs         float64   `orm:"column(carbs)"`
	Fat           float64   `orm:"column(fat)"`
	CreatedAt     time.Time `orm:"column(created_at);auto_now_add"`
}

func (m *MacroEntry) TableName() string {
	return "macro_entries"
}

func init() {
	orm.RegisterModel(&MacroEntry{})
}

func CreateMacroEntry(userID int64, date time.Time, foodName string, servingWeight float64, servingUnit string, protein, carbs, fat float64) (*MacroEntry, error) {
	o := orm.NewOrm()
	e := &MacroEntry{
		UserID:        userID,
		Date:          date,
		FoodName:      foodName,
		ServingWeight: servingWeight,
		ServingUnit:   servingUnit,
		Protein:       protein,
		Carbs:         carbs,
		Fat:           fat,
	}
	_, err := o.Insert(e)
	return e, err
}

func GetMacroEntriesByUser(userID int64) ([]*MacroEntry, error) {
	o := orm.NewOrm()
	var entries []*MacroEntry
	_, err := o.QueryTable(&MacroEntry{}).Filter("UserID", userID).OrderBy("-Date", "-CreatedAt").All(&entries)
	return entries, err
}

func GetMacroEntryByID(id, userID int64) (*MacroEntry, error) {
	o := orm.NewOrm()
	e := &MacroEntry{ID: id}
	if err := o.Read(e); err != nil || e.UserID != userID {
		return nil, orm.ErrNoRows
	}
	return e, nil
}

func UpdateMacroEntry(id, userID int64, foodName string, servingWeight float64, servingUnit string, protein, carbs, fat float64) (*MacroEntry, error) {
	e, err := GetMacroEntryByID(id, userID)
	if err != nil {
		return nil, err
	}
	e.FoodName = foodName
	e.ServingWeight = servingWeight
	e.ServingUnit = servingUnit
	e.Protein = protein
	e.Carbs = carbs
	e.Fat = fat
	o := orm.NewOrm()
	_, err = o.Update(e, "food_name", "serving_weight", "serving_unit", "protein", "carbs", "fat")
	return e, err
}

func DeleteMacroEntry(id, userID int64) error {
	e, err := GetMacroEntryByID(id, userID)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	_, err = o.Delete(e)
	return err
}
