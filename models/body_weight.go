package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type BodyWeight struct {
	ID         int64     `orm:"column(id);auto;pk"`
	UserID     int64     `orm:"column(user_id)"`
	Date       time.Time `orm:"column(date);type(date)"`
	Weight     float64   `orm:"column(weight)"`
	WeightUnit string    `orm:"column(weight_unit)"`
}

func (b *BodyWeight) TableName() string {
	return "body_weights"
}

func init() {
	orm.RegisterModel(&BodyWeight{})
}

func CreateBodyWeight(userID int64, date time.Time, weight float64, weightUnit string) (*BodyWeight, error) {
	o := orm.NewOrm()
	bw := &BodyWeight{
		UserID:     userID,
		Date:       date,
		Weight:     weight,
		WeightUnit: weightUnit,
	}
	_, err := o.Insert(bw)
	return bw, err
}

func GetBodyWeightsByUser(userID int64) ([]*BodyWeight, error) {
	o := orm.NewOrm()
	var entries []*BodyWeight
	_, err := o.QueryTable(&BodyWeight{}).Filter("UserID", userID).OrderBy("-Date").All(&entries)
	return entries, err
}

func GetBodyWeightByID(id, userID int64) (*BodyWeight, error) {
	o := orm.NewOrm()
	bw := &BodyWeight{ID: id}
	if err := o.Read(bw); err != nil || bw.UserID != userID {
		return nil, orm.ErrNoRows
	}
	return bw, nil
}

func UpdateBodyWeight(id, userID int64, weight float64, weightUnit string) (*BodyWeight, error) {
	bw, err := GetBodyWeightByID(id, userID)
	if err != nil {
		return nil, err
	}
	bw.Weight = weight
	bw.WeightUnit = weightUnit
	o := orm.NewOrm()
	_, err = o.Update(bw, "weight", "weight_unit")
	return bw, err
}

func DeleteBodyWeight(id, userID int64) error {
	bw, err := GetBodyWeightByID(id, userID)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	_, err = o.Delete(bw)
	return err
}
