package models

import (
	"github.com/beego/beego/v2/client/orm"
)

type SessionExercise struct {
	ID           int64   `orm:"column(id);auto;pk"`
	SessionID    int64   `orm:"column(session_id)"`
	Name         string  `orm:"column(name)"`
	IsBodyweight bool    `orm:"column(is_bodyweight)"`
	GoalWeight   float64 `orm:"column(goal_weight)"`
	WeightUnit   string  `orm:"column(weight_unit)"`
	GoalReps     int     `orm:"column(goal_reps)"`
	SortOrder    int     `orm:"column(sort_order)"`
}

func (s *SessionExercise) TableName() string {
	return "session_exercises"
}

type SessionSet struct {
	ID                int64   `orm:"column(id);auto;pk"`
	SessionExerciseID int64   `orm:"column(session_exercise_id)"`
	SetNumber         int     `orm:"column(set_number)"`
	ActualWeight      float64 `orm:"column(actual_weight)"`
	WeightUnit        string  `orm:"column(weight_unit)"`
	ActualReps        int     `orm:"column(actual_reps)"`
}

func (s *SessionSet) TableName() string {
	return "session_sets"
}

// SessionExerciseView bundles an exercise with its logged sets for display.
type SessionExerciseView struct {
	Exercise *SessionExercise
	Sets     []*SessionSet
}

func init() {
	orm.RegisterModel(&SessionExercise{}, &SessionSet{})
}

func CreateSessionExercise(sessionID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, goalReps int) (*SessionExercise, error) {
	o := orm.NewOrm()
	n, _ := o.QueryTable(&SessionExercise{}).Filter("SessionID", sessionID).Count()
	e := &SessionExercise{
		SessionID:    sessionID,
		Name:         name,
		IsBodyweight: isBodyweight,
		GoalWeight:   goalWeight,
		WeightUnit:   weightUnit,
		GoalReps:     goalReps,
		SortOrder:    int(n),
	}
	_, err := o.Insert(e)
	return e, err
}

func GetSessionExercisesWithSets(sessionID int64) ([]*SessionExerciseView, error) {
	o := orm.NewOrm()
	var exercises []*SessionExercise
	_, err := o.QueryTable(&SessionExercise{}).Filter("SessionID", sessionID).OrderBy("SortOrder").All(&exercises)
	if err != nil {
		return nil, err
	}
	views := make([]*SessionExerciseView, len(exercises))
	for i, ex := range exercises {
		var sets []*SessionSet
		o.QueryTable(&SessionSet{}).Filter("SessionExerciseID", ex.ID).OrderBy("SetNumber").All(&sets)
		views[i] = &SessionExerciseView{Exercise: ex, Sets: sets}
	}
	return views, nil
}

func GetSessionExerciseByID(exerciseID int64) (*SessionExercise, error) {
	o := orm.NewOrm()
	e := &SessionExercise{ID: exerciseID}
	if err := o.Read(e); err != nil {
		return nil, err
	}
	return e, nil
}

func LogSessionSet(exerciseID int64, setNumber int, actualWeight float64, weightUnit string, actualReps int) (*SessionSet, error) {
	o := orm.NewOrm()
	s := &SessionSet{
		SessionExerciseID: exerciseID,
		SetNumber:         setNumber,
		ActualWeight:      actualWeight,
		WeightUnit:        weightUnit,
		ActualReps:        actualReps,
	}
	_, err := o.Insert(s)
	return s, err
}

func CountSetsByExercise(exerciseID int64) (int, error) {
	o := orm.NewOrm()
	n, err := o.QueryTable(&SessionSet{}).Filter("SessionExerciseID", exerciseID).Count()
	return int(n), err
}
