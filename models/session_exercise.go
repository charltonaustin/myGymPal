package models

import (
	"errors"
	"strings"

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
	Block        string  `orm:"column(block)"`
	SortOrder    int     `orm:"column(sort_order)"`
	IsTimeBased  bool    `orm:"column(is_time_based)"`
	GoalSeconds  int     `orm:"column(goal_seconds)"`
	// LinkedToNext means "do not rest after me — go straight to the next exercise".
	// It is a property of this exercise alone, not of a pair, so reordering and
	// deleting stay safe. A stale true on an exercise that ends up last in its
	// block is ignored at render time; never read it directly for display.
	LinkedToNext bool `orm:"column(linked_to_next)"`
	// CircuitID is the circuit this exercise runs inside, or nil when it is a
	// loose exercise. Membership is this field and only this field — never a
	// non-zero WorkSeconds, which a loose exercise is permitted to carry.
	CircuitID *int64 `orm:"column(circuit_id);null"`
	// WorkSeconds is how long one round of this exercise lasts. Only meaningful
	// inside a circuit; 0 otherwise.
	WorkSeconds int `orm:"column(work_seconds)"`
}

func (s *SessionExercise) TableName() string {
	return "session_exercises"
}

// SessionCircuit is a named, ordered group of exercises inside one session. The
// group runs Rounds times and TransitionSeconds separates one exercise from the
// next. It is copied from a template circuit, never linked to one: editing the
// template afterwards must not rewrite a workout already performed.
type SessionCircuit struct {
	ID                int64  `orm:"column(id);auto;pk"`
	SessionID         int64  `orm:"column(session_id)"`
	Name              string `orm:"column(name)"`
	Rounds            int    `orm:"column(rounds)"`
	TransitionSeconds int    `orm:"column(transition_seconds)"`
	SortOrder         int    `orm:"column(sort_order)"`
}

func (c *SessionCircuit) TableName() string {
	return "session_circuits"
}

// SessionExerciseInput is everything needed to create one exercise in a session.
// It replaces a nine-parameter constructor: adding circuit membership and work
// seconds would have made eleven positional arguments, almost all of them int
// and bool, where a transposed pair still compiles and still passes the tests.
type SessionExerciseInput struct {
	SessionID    int64
	Name         string
	IsBodyweight bool
	GoalWeight   float64
	WeightUnit   string
	GoalReps     int
	Block        string
	IsTimeBased  bool
	GoalSeconds  int
	// CircuitID is the circuit this exercise belongs to, nil when it is loose. A
	// pointer rather than an index with a sentinel, so a struct literal that
	// omits the field means "loose" instead of "the first circuit".
	CircuitID   *int64
	WorkSeconds int
}

// SessionCircuitInput is one circuit being copied into a session.
type SessionCircuitInput struct {
	// TemplateCircuitID is the circuit this one is copied from. It is the key
	// the copied exercises' CircuitID values are matched against, and it is
	// deliberately not stored: a session circuit is a copy, not a reference.
	TemplateCircuitID int64
	Name              string
	Rounds            int
	TransitionSeconds int
	SortOrder         int
}

type SessionSet struct {
	ID                int64   `orm:"column(id);auto;pk"`
	SessionExerciseID int64   `orm:"column(session_exercise_id)"`
	SetNumber         int     `orm:"column(set_number)"`
	ActualWeight      float64 `orm:"column(actual_weight)"`
	WeightUnit        string  `orm:"column(weight_unit)"`
	ActualReps        int     `orm:"column(actual_reps)"`
	ActualSeconds     int     `orm:"column(actual_seconds)"`
	ActivityType      string  `orm:"column(activity_type)"`
}

func (s *SessionSet) TableName() string {
	return "session_sets"
}

func (s *SessionSet) Hours() int   { return s.ActualSeconds / 3600 }
func (s *SessionSet) Minutes() int { return (s.ActualSeconds % 3600) / 60 }
func (s *SessionSet) Secs() int    { return s.ActualSeconds % 60 }

// SessionExerciseView bundles an exercise with its logged sets and cardio logs for display.
type SessionExerciseView struct {
	Exercise   *SessionExercise
	Sets       []*SessionSet
	CardioLogs []*CardioLog
	HitMax     bool // true if the user hit max reps at goal weight for all required sets in the previous session
	BelowGoal  bool // true if the user logged any set below goal weight in the previous session
	GoalRepMin int  // from exercise library; overrides phase rep range for bodyweight exercises when > 0
	GoalRepMax int  // from exercise library; overrides phase rep range for bodyweight exercises when > 0

	// SupersetLinked is the effective link, computed by the controller: the raw
	// LinkedToNext, but only when there is a next exercise in the same block and
	// the run is still under the four-member cap.
	SupersetLinked bool
	// SupersetLabel is "A1", "A2", … for a member of a superset run, or "" for a
	// solo exercise. Assigned per block, in sort order.
	SupersetLabel string
}

// LastSet returns the most recently logged set, or nil if none exist.
func (v *SessionExerciseView) LastSet() *SessionSet {
	if len(v.Sets) == 0 {
		return nil
	}
	return v.Sets[len(v.Sets)-1]
}

func init() {
	orm.RegisterModel(&SessionExercise{}, &SessionSet{}, &SessionCircuit{})
}

// newSessionExercise builds the row for one input at a known sort order. It is
// shared by the single-exercise create and the template copy so that a field one
// of them writes cannot go missing from the other.
func newSessionExercise(in SessionExerciseInput, sortOrder int) *SessionExercise {
	block := in.Block
	if block == "" {
		block = "main"
	}
	return &SessionExercise{
		SessionID:    in.SessionID,
		Name:         strings.ToLower(strings.TrimSpace(in.Name)),
		IsBodyweight: in.IsBodyweight,
		GoalWeight:   in.GoalWeight,
		WeightUnit:   in.WeightUnit,
		GoalReps:     in.GoalReps,
		Block:        block,
		SortOrder:    sortOrder,
		IsTimeBased:  in.IsTimeBased,
		GoalSeconds:  in.GoalSeconds,
		CircuitID:    in.CircuitID,
		WorkSeconds:  ValidSeconds(in.WorkSeconds),
	}
}

func CreateSessionExercise(in SessionExerciseInput) (*SessionExercise, error) {
	o := orm.NewOrm()
	n, _ := o.QueryTable(&SessionExercise{}).Filter("SessionID", in.SessionID).Count()
	e := newSessionExercise(in, int(n))
	_, err := o.Insert(e)
	return e, err
}

// CreateSessionBody copies a template's circuits and exercises into a session in
// one transaction. Every circuit is inserted first, then every exercise, with
// each exercise's CircuitID — which arrives holding the *template* circuit id —
// rewritten to the id of the session circuit copied from it.
//
// One transaction is the point: the previous code inserted one exercise per call
// with no transaction, so a failure partway through left a session holding half
// a workout with no way to tell. An exercise naming a circuit that was not
// submitted alongside it becomes loose rather than dangling, so one bad row
// cannot roll back the whole workout on a foreign key violation.
func CreateSessionBody(sessionID int64, circuits []SessionCircuitInput, exercises []SessionExerciseInput) error {
	o := orm.NewOrm()
	base, _ := o.QueryTable(&SessionExercise{}).Filter("SessionID", sessionID).Count()

	tx, err := o.Begin()
	if err != nil {
		return err
	}

	// Template circuit id → the id of the session circuit copied from it.
	idByTemplateID := make(map[int64]int64, len(circuits))
	for _, ci := range circuits {
		c := &SessionCircuit{
			SessionID:         sessionID,
			Name:              strings.TrimSpace(ci.Name),
			Rounds:            ValidRounds(ci.Rounds),
			TransitionSeconds: ValidSeconds(ci.TransitionSeconds),
			SortOrder:         ci.SortOrder,
		}
		if _, err := tx.Insert(c); err != nil {
			tx.Rollback()
			return err
		}
		idByTemplateID[ci.TemplateCircuitID] = c.ID
	}

	for i, in := range exercises {
		in.SessionID = sessionID
		if in.CircuitID != nil {
			if newID, ok := idByTemplateID[*in.CircuitID]; ok {
				in.CircuitID = &newID
			} else {
				in.CircuitID = nil
			}
		}
		if _, err := tx.Insert(newSessionExercise(in, int(base)+i)); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// GetSessionCircuits reads a session's circuits in display order. Like the
// exercise reads it is scoped by session alone; the caller owns the check that
// the session belongs to the user.
func GetSessionCircuits(sessionID int64) ([]*SessionCircuit, error) {
	o := orm.NewOrm()
	var circuits []*SessionCircuit
	_, err := o.QueryTable(&SessionCircuit{}).Filter("SessionID", sessionID).OrderBy("SortOrder").All(&circuits)
	if err != nil {
		return nil, err
	}
	return circuits, nil
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
		var cardioLogs []*CardioLog
		if ex.Block == "cardio" {
			o.QueryTable(&CardioLog{}).Filter("SessionExerciseID", ex.ID).OrderBy("CreatedAt").All(&cardioLogs)
		}
		views[i] = &SessionExerciseView{Exercise: ex, Sets: sets, CardioLogs: cardioLogs}
	}
	return views, nil
}

func DeleteSessionExercise(exerciseID int64) error {
	o := orm.NewOrm()
	_, err := o.Delete(&SessionExercise{ID: exerciseID})
	return err
}

func GetSessionExerciseByID(exerciseID int64) (*SessionExercise, error) {
	o := orm.NewOrm()
	e := &SessionExercise{ID: exerciseID}
	if err := o.Read(e); err != nil {
		return nil, err
	}
	return e, nil
}

func UpdateSessionExerciseName(id int64, name string) error {
	o := orm.NewOrm()
	_, err := o.Raw("UPDATE session_exercises SET name = ? WHERE id = ?", name, id).Exec()
	return err
}

func UpdateSessionExerciseLink(id int64, linked bool) error {
	o := orm.NewOrm()
	_, err := o.Raw("UPDATE session_exercises SET linked_to_next = ? WHERE id = ?", linked, id).Exec()
	return err
}

func LogSessionSet(exerciseID int64, setNumber int, actualWeight float64, weightUnit string, actualReps int, actualSeconds int, activityType string) (*SessionSet, error) {
	o := orm.NewOrm()
	s := &SessionSet{
		SessionExerciseID: exerciseID,
		SetNumber:         setNumber,
		ActualWeight:      actualWeight,
		WeightUnit:        weightUnit,
		ActualReps:        actualReps,
		ActualSeconds:     actualSeconds,
		ActivityType:      activityType,
	}
	_, err := o.Insert(s)
	return s, err
}

func CountSetsByExercise(exerciseID int64) (int, error) {
	o := orm.NewOrm()
	n, err := o.QueryTable(&SessionSet{}).Filter("SessionExerciseID", exerciseID).Count()
	return int(n), err
}

func DeleteSessionSet(setID int64) error {
	o := orm.NewOrm()
	s := &SessionSet{ID: setID}
	if err := o.Read(s); err != nil {
		return errors.New("not found")
	}
	tx, err := o.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Raw("DELETE FROM session_sets WHERE id = ?", setID).Exec(); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Raw("UPDATE session_sets SET set_number = set_number - 1 WHERE session_exercise_id = ? AND set_number > ?", s.SessionExerciseID, s.SetNumber).Exec(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func UpdateSessionExerciseSortOrders(sessionID int64, ids []int64) error {
	o := orm.NewOrm()
	for i, id := range ids {
		if _, err := o.Raw("UPDATE session_exercises SET sort_order = ? WHERE id = ? AND session_id = ?", i, id, sessionID).Exec(); err != nil {
			return err
		}
	}
	return nil
}
