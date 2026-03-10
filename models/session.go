package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Session struct {
	ID            int64     `orm:"column(id);auto;pk"`
	ProgramID     int64     `orm:"column(program_id)"`
	UserID        int64     `orm:"column(user_id)"`
	PhaseNumber   int       `orm:"column(phase_number)"`
	WeekNumber    int       `orm:"column(week_number)"`
	WorkoutNumber int       `orm:"column(workout_number)"`
	IsDeload      bool      `orm:"column(is_deload)"`
	Date          time.Time `orm:"column(date);type(date)"`
	CreatedAt     time.Time `orm:"column(created_at);auto_now_add"`
	UpdatedAt     time.Time `orm:"column(updated_at);auto_now"`
}

func (s *Session) TableName() string {
	return "sessions"
}

func init() {
	orm.RegisterModel(&Session{})
}

// CalculatePhaseAndWeek returns the 1-indexed phase and week for the given
// reference time, given a program start date and weeks-per-phase setting.
// The last week of each phase is flagged as a deload week.
func CalculatePhaseAndWeek(now, startDate time.Time, weeksPerPhase int) (phase, week int, isDeload bool) {
	if weeksPerPhase <= 0 {
		weeksPerPhase = 8
	}
	daysSinceStart := int(now.Sub(startDate).Hours() / 24)
	if daysSinceStart < 0 {
		daysSinceStart = 0
	}
	totalWeeks := daysSinceStart / 7
	phase = totalWeeks/weeksPerPhase + 1
	week = totalWeeks%weeksPerPhase + 1
	isDeload = week == weeksPerPhase
	return
}

func CreateSession(programID, userID int64, phaseNumber, weekNumber, workoutNumber int, isDeload bool, date time.Time) (*Session, error) {
	o := orm.NewOrm()
	s := &Session{
		ProgramID:     programID,
		UserID:        userID,
		PhaseNumber:   phaseNumber,
		WeekNumber:    weekNumber,
		WorkoutNumber: workoutNumber,
		IsDeload:      isDeload,
		Date:          date,
	}
	_, err := o.Insert(s)
	return s, err
}

func CountSessionsByProgram(programID int64) (int, error) {
	o := orm.NewOrm()
	n, err := o.QueryTable(&Session{}).Filter("ProgramID", programID).Count()
	return int(n), err
}

func GetSessionsByProgram(programID int64) ([]*Session, error) {
	o := orm.NewOrm()
	var sessions []*Session
	_, err := o.QueryTable(&Session{}).Filter("ProgramID", programID).OrderBy("-Date", "-ID").All(&sessions)
	return sessions, err
}

func DeleteSession(id, userID int64) error {
	o := orm.NewOrm()
	s := &Session{ID: id}
	if err := o.Read(s); err != nil {
		return errors.New("not found")
	}
	if s.UserID != userID {
		return errors.New("not found")
	}
	_, err := o.Delete(s)
	return err
}

func GetSessionByID(id, userID int64) (*Session, error) {
	o := orm.NewOrm()
	s := &Session{ID: id}
	if err := o.Read(s); err != nil {
		return nil, errors.New("not found")
	}
	if s.UserID != userID {
		return nil, errors.New("not found")
	}
	return s, nil
}
