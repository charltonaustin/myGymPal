package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type CardioLog struct {
	ID                int64     `orm:"column(id);auto;pk"`
	SessionExerciseID int64     `orm:"column(session_exercise_id)"`
	CardioType        string    `orm:"column(cardio_type)"`
	GoalDuration      int       `orm:"column(goal_duration)"`
	ActualDuration    int       `orm:"column(actual_duration)"`
	CreatedAt         time.Time `orm:"column(created_at);auto_now_add"`
}

func (c *CardioLog) TableName() string {
	return "cardio_logs"
}

func init() {
	orm.RegisterModel(&CardioLog{})
}

func LogCardioEntry(sessionExerciseID int64, cardioType string, goalDuration, actualDuration int) (*CardioLog, error) {
	o := orm.NewOrm()
	log := &CardioLog{
		SessionExerciseID: sessionExerciseID,
		CardioType:        cardioType,
		GoalDuration:      goalDuration,
		ActualDuration:    actualDuration,
	}
	_, err := o.Insert(log)
	return log, err
}

func GetCardioLogsByExercise(sessionExerciseID int64) ([]*CardioLog, error) {
	o := orm.NewOrm()
	var logs []*CardioLog
	_, err := o.QueryTable(&CardioLog{}).Filter("SessionExerciseID", sessionExerciseID).OrderBy("CreatedAt").All(&logs)
	return logs, err
}

func DeleteCardioLog(id int64) error {
	o := orm.NewOrm()
	_, err := o.Delete(&CardioLog{ID: id})
	return err
}
