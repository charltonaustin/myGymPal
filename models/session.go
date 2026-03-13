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

// CalculateNextSession returns the 1-indexed phase, week (within phase), workout number
// (within week), and deload flag for the next session to be created, based on how many
// sessions already exist for the program.
func CalculateNextSession(sessionCount, weeksPerPhase, workoutsPerWeek int) (phase, week, workoutNumber int, isDeload bool) {
	if weeksPerPhase <= 0 {
		weeksPerPhase = 8
	}
	if workoutsPerWeek <= 0 {
		workoutsPerWeek = 4
	}
	totalWeeks := sessionCount / workoutsPerWeek
	phase = totalWeeks/weeksPerPhase + 1
	week = totalWeeks%weeksPerPhase + 1
	workoutNumber = sessionCount%workoutsPerWeek + 1
	isDeload = week == weeksPerPhase
	return
}

// CalculateNextSessionFromLast returns the next phase, week, workout number, and deload
// flag by incrementing from the most recent session's values.
func CalculateNextSessionFromLast(last *Session, weeksPerPhase, workoutsPerWeek int) (phase, week, workoutNumber int, isDeload bool) {
	if weeksPerPhase <= 0 {
		weeksPerPhase = 8
	}
	if workoutsPerWeek <= 0 {
		workoutsPerWeek = 4
	}
	phase = last.PhaseNumber
	week = last.WeekNumber
	workoutNumber = last.WorkoutNumber + 1
	if workoutNumber > workoutsPerWeek {
		workoutNumber = 1
		week++
		if week > weeksPerPhase {
			week = 1
			phase++
		}
	}
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

func GetLatestSessionByProgram(programID int64) (*Session, error) {
	o := orm.NewOrm()
	var s Session
	err := o.QueryTable(&Session{}).Filter("ProgramID", programID).OrderBy("-Date", "-ID").Limit(1).One(&s)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	return &s, err
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

// RecentSession bundles a session with its program name for dashboard display.
type RecentSession struct {
	ID            int64     `orm:"column(id)"`
	ProgramID     int64     `orm:"column(program_id)"`
	PhaseNumber   int       `orm:"column(phase_number)"`
	WeekNumber    int       `orm:"column(week_number)"`
	WorkoutNumber int       `orm:"column(workout_number)"`
	IsDeload      bool      `orm:"column(is_deload)"`
	Date          time.Time `orm:"column(date)"`
	ProgramName   string    `orm:"column(program_name)"`
}

func GetRecentSessionsByUser(userID int64, limit int) ([]*RecentSession, error) {
	o := orm.NewOrm()
	var rows []*RecentSession
	_, err := o.Raw(`
		SELECT s.*, p.name AS program_name
		FROM sessions s
		JOIN programs p ON p.id = s.program_id
		WHERE s.user_id = ?
		ORDER BY s.date DESC, s.id DESC
		LIMIT ?`, userID, limit).QueryRows(&rows)
	return rows, err
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
