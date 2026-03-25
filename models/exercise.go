package models

import (
	"errors"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Exercise struct {
	ID           int64     `orm:"column(id);auto;pk"`
	UserID       int64     `orm:"column(user_id)"`
	Name         string    `orm:"column(name)"`
	IsBodyweight bool      `orm:"column(is_bodyweight)"`
	GoalWeight   float64   `orm:"column(goal_weight)"`
	WeightUnit   string    `orm:"column(weight_unit)"`
	IsTimeBased  bool      `orm:"column(is_time_based)"`
	GoalSeconds  int       `orm:"column(goal_seconds)"`
	GoalRepMin   int       `orm:"column(goal_rep_min)"`
	GoalRepMax   int       `orm:"column(goal_rep_max)"`
	DefaultBlock string    `orm:"column(default_block)"`
	CreatedAt    time.Time `orm:"column(created_at);auto_now_add"`
	UpdatedAt    time.Time `orm:"column(updated_at);auto_now"`
}

func (e *Exercise) TableName() string { return "exercises" }

func init() {
	orm.RegisterModel(&Exercise{})
}

func validDefaultBlock(b string) string {
	switch b {
	case "abs", "cardio", "stretch":
		return b
	default:
		return "main"
	}
}

func CreateExercise(userID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, isTimeBased bool, goalSeconds int, goalRepMin int, goalRepMax int, defaultBlock string) (*Exercise, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return nil, errors.New("exercise name is required")
	}
	o := orm.NewOrm()
	ex := &Exercise{
		UserID:       userID,
		Name:         name,
		IsBodyweight: isBodyweight,
		GoalWeight:   goalWeight,
		WeightUnit:   weightUnit,
		IsTimeBased:  isTimeBased,
		GoalSeconds:  goalSeconds,
		GoalRepMin:   goalRepMin,
		GoalRepMax:   goalRepMax,
		DefaultBlock: validDefaultBlock(defaultBlock),
	}
	if _, err := o.Insert(ex); err != nil {
		return nil, err
	}
	return ex, nil
}

func GetExercisesByUser(userID int64) ([]*Exercise, error) {
	o := orm.NewOrm()
	var exercises []*Exercise
	_, err := o.QueryTable(&Exercise{}).Filter("UserID", userID).OrderBy("Name").All(&exercises)
	return exercises, err
}

func GetExerciseByID(id, userID int64) (*Exercise, error) {
	o := orm.NewOrm()
	ex := &Exercise{ID: id}
	if err := o.Read(ex); err != nil {
		return nil, errors.New("not found")
	}
	if ex.UserID != userID {
		return nil, errors.New("not found")
	}
	return ex, nil
}

func GetExerciseByName(userID int64, name string) (*Exercise, error) {
	o := orm.NewOrm()
	var exercises []*Exercise
	_, err := o.Raw(
		"SELECT * FROM exercises WHERE user_id = ? AND LOWER(TRIM(name)) = LOWER(TRIM(?)) LIMIT 1",
		userID, name,
	).QueryRows(&exercises)
	if err != nil || len(exercises) == 0 {
		return nil, errors.New("not found")
	}
	return exercises[0], nil
}

func UpdateExercise(id, userID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, isTimeBased bool, goalSeconds int, goalRepMin int, goalRepMax int, defaultBlock string) (*Exercise, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return nil, errors.New("exercise name is required")
	}
	ex, err := GetExerciseByID(id, userID)
	if err != nil {
		return nil, err
	}
	ex.Name = name
	ex.IsBodyweight = isBodyweight
	ex.GoalWeight = goalWeight
	ex.WeightUnit = weightUnit
	ex.IsTimeBased = isTimeBased
	ex.GoalSeconds = goalSeconds
	ex.GoalRepMin = goalRepMin
	ex.GoalRepMax = goalRepMax
	ex.DefaultBlock = validDefaultBlock(defaultBlock)
	o := orm.NewOrm()
	if _, err := o.Update(ex); err != nil {
		return nil, err
	}
	return ex, nil
}

func UpdateExerciseGoalWeight(id int64, goalWeight float64) error {
	o := orm.NewOrm()
	ex := &Exercise{ID: id}
	if err := o.Read(ex); err != nil {
		return errors.New("not found")
	}
	ex.GoalWeight = goalWeight
	_, err := o.Update(ex, "GoalWeight", "UpdatedAt")
	return err
}

func DeleteExercise(id, userID int64) error {
	ex, err := GetExerciseByID(id, userID)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	_, err = o.Delete(ex)
	return err
}
