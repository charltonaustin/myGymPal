package models

import (
	"errors"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Exercise is a view type combining a global exercise name with a user's personal goals.
// It is not mapped to a single table; all queries use raw SQL with a JOIN.
type Exercise struct {
	ID           int64     `orm:"column(id)"`
	Name         string    `orm:"column(name)"`
	IsBodyweight bool      `orm:"column(is_bodyweight)"`
	GoalWeight   float64   `orm:"column(goal_weight)"`
	WeightUnit   string    `orm:"column(weight_unit)"`
	IsTimeBased  bool      `orm:"column(is_time_based)"`
	GoalSeconds  int       `orm:"column(goal_seconds)"`
	GoalRepMin   int       `orm:"column(goal_rep_min)"`
	GoalRepMax   int       `orm:"column(goal_rep_max)"`
	DefaultBlock string    `orm:"column(default_block)"`
	CreatedAt    time.Time `orm:"column(created_at)"`
	UpdatedAt    time.Time `orm:"column(updated_at)"`
}

// exerciseGlobal is the ORM model for the global exercises table.
type exerciseGlobal struct {
	ID        int64     `orm:"column(id);auto;pk"`
	Name      string    `orm:"column(name)"`
	CreatedAt time.Time `orm:"column(created_at);auto_now_add"`
}

func (e *exerciseGlobal) TableName() string { return "exercises" }

// userExerciseGoal is the ORM model for per-user exercise configuration.
type userExerciseGoal struct {
	ID           int64     `orm:"column(id);auto;pk"`
	UserID       int64     `orm:"column(user_id)"`
	ExerciseID   int64     `orm:"column(exercise_id)"`
	IsBodyweight bool      `orm:"column(is_bodyweight)"`
	IsTimeBased  bool      `orm:"column(is_time_based)"`
	GoalWeight   float64   `orm:"column(goal_weight)"`
	WeightUnit   string    `orm:"column(weight_unit)"`
	GoalSeconds  int       `orm:"column(goal_seconds)"`
	GoalRepMin   int       `orm:"column(goal_rep_min)"`
	GoalRepMax   int       `orm:"column(goal_rep_max)"`
	DefaultBlock string    `orm:"column(default_block)"`
	CreatedAt    time.Time `orm:"column(created_at);auto_now_add"`
	UpdatedAt    time.Time `orm:"column(updated_at);auto_now"`
}

func (e *userExerciseGoal) TableName() string { return "user_exercise_goals" }

func init() {
	orm.RegisterModel(&exerciseGlobal{}, &userExerciseGoal{})
}

// ValidBlock coerces an untrusted block name to one of the four known sections,
// falling back to "main". Callers pass raw form input straight into it.
func ValidBlock(b string) string {
	switch b {
	case "abs", "cardio", "stretch":
		return b
	default:
		return "main"
	}
}

// exerciseJoinSQL selects all Exercise fields via LEFT JOIN with user goals.
// The caller appends a WHERE clause and passes (userID, ...) as args.
const exerciseJoinSQL = `
SELECT
    e.id,
    e.name,
    COALESCE(g.is_bodyweight, false) AS is_bodyweight,
    COALESCE(g.is_time_based,  false) AS is_time_based,
    COALESCE(g.goal_weight,    0)     AS goal_weight,
    COALESCE(g.weight_unit,   'lb')   AS weight_unit,
    COALESCE(g.goal_seconds,   0)     AS goal_seconds,
    COALESCE(g.goal_rep_min,   0)     AS goal_rep_min,
    COALESCE(g.goal_rep_max,   0)     AS goal_rep_max,
    COALESCE(g.default_block, 'main') AS default_block,
    e.created_at,
    COALESCE(g.updated_at, e.created_at) AS updated_at
FROM exercises e
LEFT JOIN user_exercise_goals g ON g.exercise_id = e.id AND g.user_id = ?
`

// exerciseInnerJoinSQL selects only exercises the user has configured (INNER JOIN).
const exerciseInnerJoinSQL = `
SELECT
    e.id,
    e.name,
    g.is_bodyweight,
    g.is_time_based,
    g.goal_weight,
    g.weight_unit,
    g.goal_seconds,
    g.goal_rep_min,
    g.goal_rep_max,
    g.default_block,
    e.created_at,
    g.updated_at
FROM exercises e
INNER JOIN user_exercise_goals g ON g.exercise_id = e.id AND g.user_id = ?
`

func upsertUserGoal(userID, exerciseID int64, isBodyweight bool, goalWeight float64, weightUnit string, isTimeBased bool, goalSeconds, goalRepMin, goalRepMax int, defaultBlock string) error {
	o := orm.NewOrm()
	_, err := o.Raw(`
		INSERT INTO user_exercise_goals
		    (user_id, exercise_id, is_bodyweight, is_time_based, goal_weight, weight_unit, goal_seconds, goal_rep_min, goal_rep_max, default_block)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (user_id, exercise_id) DO UPDATE SET
		    is_bodyweight = EXCLUDED.is_bodyweight,
		    is_time_based = EXCLUDED.is_time_based,
		    goal_weight   = EXCLUDED.goal_weight,
		    weight_unit   = EXCLUDED.weight_unit,
		    goal_seconds  = EXCLUDED.goal_seconds,
		    goal_rep_min  = EXCLUDED.goal_rep_min,
		    goal_rep_max  = EXCLUDED.goal_rep_max,
		    default_block = EXCLUDED.default_block,
		    updated_at    = NOW()
	`, userID, exerciseID, isBodyweight, isTimeBased, goalWeight, weightUnit, goalSeconds, goalRepMin, goalRepMax, defaultBlock).Exec()
	return err
}

func CreateExercise(userID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, isTimeBased bool, goalSeconds, goalRepMin, goalRepMax int, defaultBlock string) (*Exercise, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return nil, errors.New("exercise name is required")
	}

	// Check if user already has this exercise configured.
	if _, err := GetExerciseByName(userID, name); err == nil {
		return nil, errors.New("exercise already exists in your library")
	}

	o := orm.NewOrm()
	var exerciseID int64
	if err := o.Raw(`INSERT INTO exercises (name) VALUES (?) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id`, name).QueryRow(&exerciseID); err != nil {
		return nil, err
	}

	if err := upsertUserGoal(userID, exerciseID, isBodyweight, goalWeight, weightUnit, isTimeBased, goalSeconds, goalRepMin, goalRepMax, ValidBlock(defaultBlock)); err != nil {
		return nil, err
	}

	return GetExerciseByID(exerciseID, userID)
}

// GetExercisesAll returns all global exercises with the user's goals overlaid where configured.
// Used for exercise autocomplete (shows all global exercises).
func GetExercisesAll(userID int64) ([]*Exercise, error) {
	o := orm.NewOrm()
	var exercises []*Exercise
	_, err := o.Raw(exerciseJoinSQL+"ORDER BY e.name", userID).QueryRows(&exercises)
	return exercises, err
}

// GetGlobalExercisesNotConfigured returns names of global exercises the user has no goal for yet.
// Used to populate autocomplete on the new exercise form.
func GetGlobalExercisesNotConfigured(userID int64) ([]string, error) {
	o := orm.NewOrm()
	var names []string
	_, err := o.Raw(`
		SELECT e.name FROM exercises e
		LEFT JOIN user_exercise_goals g ON g.exercise_id = e.id AND g.user_id = ?
		WHERE g.exercise_id IS NULL
		ORDER BY e.name
	`, userID).QueryRows(&names)
	return names, err
}

// GetExercisesByUser returns only exercises the user has personally configured.
// Used for the exercise management page.
func GetExercisesByUser(userID int64) ([]*Exercise, error) {
	o := orm.NewOrm()
	var exercises []*Exercise
	_, err := o.Raw(exerciseInnerJoinSQL+"ORDER BY e.name", userID).QueryRows(&exercises)
	return exercises, err
}

func GetExerciseByID(id, userID int64) (*Exercise, error) {
	o := orm.NewOrm()
	var exercises []*Exercise
	_, err := o.Raw(exerciseInnerJoinSQL+"WHERE e.id = ?", userID, id).QueryRows(&exercises)
	if err != nil || len(exercises) == 0 {
		return nil, errors.New("not found")
	}
	return exercises[0], nil
}

func GetExerciseByName(userID int64, name string) (*Exercise, error) {
	o := orm.NewOrm()
	var exercises []*Exercise
	_, err := o.Raw(exerciseInnerJoinSQL+"WHERE LOWER(TRIM(e.name)) = LOWER(TRIM(?)) LIMIT 1", userID, name).QueryRows(&exercises)
	if err != nil || len(exercises) == 0 {
		return nil, errors.New("not found")
	}
	return exercises[0], nil
}

func UpdateExercise(id, userID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, isTimeBased bool, goalSeconds, goalRepMin, goalRepMax int, defaultBlock string) (*Exercise, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return nil, errors.New("exercise name is required")
	}

	o := orm.NewOrm()
	if _, err := o.Raw(`UPDATE exercises SET name = ? WHERE id = ?`, name, id).Exec(); err != nil {
		return nil, err
	}

	if err := upsertUserGoal(userID, id, isBodyweight, goalWeight, weightUnit, isTimeBased, goalSeconds, goalRepMin, goalRepMax, ValidBlock(defaultBlock)); err != nil {
		return nil, err
	}

	return GetExerciseByID(id, userID)
}

func UpdateExerciseGoalWeight(id, userID int64, goalWeight float64, weightUnit string) error {
	o := orm.NewOrm()
	_, err := o.Raw(`
		INSERT INTO user_exercise_goals (user_id, exercise_id, goal_weight, weight_unit, updated_at)
		VALUES (?, ?, ?, ?, NOW())
		ON CONFLICT (user_id, exercise_id) DO UPDATE SET
		    goal_weight = EXCLUDED.goal_weight,
		    weight_unit = EXCLUDED.weight_unit,
		    updated_at  = NOW()
	`, userID, id, goalWeight, weightUnit).Exec()
	return err
}

func DeleteExercise(id, userID int64) error {
	o := orm.NewOrm()
	_, err := o.Raw(`DELETE FROM user_exercise_goals WHERE exercise_id = ? AND user_id = ?`, id, userID).Exec()
	return err
}

// EnsureExerciseExists creates the global exercise entry if it doesn't exist yet,
// without creating or modifying user goals. Returns the global exercise ID.
func EnsureExerciseExists(name string) (int64, error) {
	o := orm.NewOrm()
	var id int64
	err := o.Raw(`INSERT INTO exercises (name) VALUES (?) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id`, name).QueryRow(&id)
	return id, err
}
