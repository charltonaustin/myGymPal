package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Template struct {
	ID        int64     `orm:"column(id);auto;pk"`
	Name      string    `orm:"column(name)"`
	Focus     string    `orm:"column(focus)"`
	CreatedAt time.Time `orm:"column(created_at);auto_now_add"`
	UpdatedAt time.Time `orm:"column(updated_at);auto_now"`
}

func (t *Template) TableName() string {
	return "templates"
}

// TemplateCircuit is a named, ordered group of exercises within a template. The
// group runs Rounds times, and TransitionSeconds of rest separates one exercise
// from the next.
type TemplateCircuit struct {
	ID                int64  `orm:"column(id);auto;pk"`
	TemplateID        int64  `orm:"column(template_id)"`
	Name              string `orm:"column(name)"`
	Rounds            int    `orm:"column(rounds)"`
	TransitionSeconds int    `orm:"column(transition_seconds)"`
	SortOrder         int    `orm:"column(sort_order)"`
}

func (c *TemplateCircuit) TableName() string {
	return "template_circuits"
}

type TemplateExercise struct {
	ID           int64  `orm:"column(id);auto;pk"`
	TemplateID   int64  `orm:"column(template_id)"`
	Name         string `orm:"column(name)"`
	IsBodyweight bool   `orm:"column(is_bodyweight)"`
	IsTimeBased  bool   `orm:"column(is_time_based)"`
	Block        string `orm:"column(block)"`
	SortOrder    int    `orm:"column(sort_order)"`
	// CircuitID is nil for a normal exercise and set for one inside a circuit.
	CircuitID *int64 `orm:"column(circuit_id);null"`
	// WorkSeconds is how long the exercise is worked for. Only meaningful when
	// the exercise is in a circuit; 0 otherwise.
	WorkSeconds int `orm:"column(work_seconds)"`
}

func (e *TemplateExercise) TableName() string {
	return "template_exercises"
}

type TemplateCircuitInput struct {
	Name              string
	Rounds            int
	TransitionSeconds int
	SortOrder         int
}

type TemplateExerciseInput struct {
	Name         string
	IsBodyweight bool
	IsTimeBased  bool
	Block        string
	SortOrder    int
	// CircuitIndex refers to a position in the circuits slice passed alongside
	// this input, because a circuit and the exercises inside it are created in
	// the same submit and the circuit has no id yet. NoCircuit means the
	// exercise is not in a circuit.
	CircuitIndex int
	WorkSeconds  int
}

// NoCircuit is the CircuitIndex of an exercise that is not part of a circuit.
const NoCircuit = -1

func init() {
	orm.RegisterModel(&Template{}, &TemplateCircuit{}, &TemplateExercise{})
}

// ValidRounds coerces an untrusted round count to at least one, mirroring the
// rounds >= 1 check on the table.
func ValidRounds(n int) int {
	if n < 1 {
		return 1
	}
	return n
}

// ValidSeconds coerces an untrusted duration to a non-negative value, mirroring
// the >= 0 checks on transition_seconds and work_seconds.
func ValidSeconds(n int) int {
	if n < 0 {
		return 0
	}
	return n
}

// normalizeTemplateInput lower-cases and trims exercise names in place and
// rejects anything the tables would reject, so Create and Update fail the same
// way on the same input.
func normalizeTemplateInput(name string, circuits []TemplateCircuitInput, exercises []TemplateExerciseInput) error {
	if name == "" {
		return errors.New("template name is required")
	}
	if len(exercises) == 0 {
		return errors.New("at least one exercise is required")
	}
	for i, ci := range circuits {
		circuits[i].Name = strings.TrimSpace(ci.Name)
		if circuits[i].Name == "" {
			return errors.New("circuit name is required")
		}
		if ci.Rounds < 1 {
			return fmt.Errorf("circuit %q must have at least one round", circuits[i].Name)
		}
		if ci.TransitionSeconds < 0 {
			return fmt.Errorf("circuit %q cannot have a negative transition", circuits[i].Name)
		}
	}
	for i, ex := range exercises {
		exercises[i].Name = strings.ToLower(strings.TrimSpace(ex.Name))
		if exercises[i].Name == "" {
			return errors.New("exercise name is required")
		}
		if ex.WorkSeconds < 0 {
			return fmt.Errorf("exercise %q cannot have negative work seconds", exercises[i].Name)
		}
		// An index pointing at a circuit that was never submitted is a bug in the
		// caller. Storing NULL instead would silently drop the exercise out of its
		// circuit, which reads as data loss to the engineer and is far harder to
		// trace back to here.
		if ex.CircuitIndex != NoCircuit && (ex.CircuitIndex < 0 || ex.CircuitIndex >= len(circuits)) {
			return fmt.Errorf("exercise %q references circuit %d, which does not exist", exercises[i].Name, ex.CircuitIndex)
		}
	}
	return nil
}

// insertTemplateBody writes a template's circuits and then its exercises,
// resolving each exercise's CircuitIndex to the id of the circuit just inserted.
// Create and Update share it so that a field added to one insert cannot go
// missing from the other — the bug that silently cleared is_time_based on every
// edit.
func insertTemplateBody(tx orm.TxOrmer, templateID int64, circuits []TemplateCircuitInput, exercises []TemplateExerciseInput) error {
	circuitIDs := make([]int64, len(circuits))
	for i, ci := range circuits {
		c := &TemplateCircuit{
			TemplateID:        templateID,
			Name:              ci.Name,
			Rounds:            ci.Rounds,
			TransitionSeconds: ci.TransitionSeconds,
			SortOrder:         ci.SortOrder,
		}
		if _, err := tx.Insert(c); err != nil {
			return err
		}
		circuitIDs[i] = c.ID
	}

	for _, ex := range exercises {
		e := &TemplateExercise{
			TemplateID:   templateID,
			Name:         ex.Name,
			IsBodyweight: ex.IsBodyweight,
			IsTimeBased:  ex.IsTimeBased,
			Block:        ex.Block,
			SortOrder:    ex.SortOrder,
			WorkSeconds:  ex.WorkSeconds,
		}
		if ex.CircuitIndex != NoCircuit {
			id := circuitIDs[ex.CircuitIndex]
			e.CircuitID = &id
		}
		if _, err := tx.Insert(e); err != nil {
			return err
		}
	}
	return nil
}

func CreateTemplate(name, focus string, circuits []TemplateCircuitInput, exercises []TemplateExerciseInput) (*Template, error) {
	if err := normalizeTemplateInput(name, circuits, exercises); err != nil {
		return nil, err
	}

	t := &Template{Name: name, Focus: focus}

	tx, err := orm.NewOrm().Begin()
	if err != nil {
		return nil, err
	}

	if _, err := tx.Insert(t); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := insertTemplateBody(tx, t.ID, circuits, exercises); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return t, nil
}

func GetAllTemplates() ([]*Template, error) {
	o := orm.NewOrm()
	var templates []*Template
	_, err := o.QueryTable(&Template{}).OrderBy("Name").All(&templates)
	return templates, err
}

func DeleteTemplate(id int64) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Template{ID: id})
	return err
}

// UpdateTemplate replaces a template's body wholesale: every circuit and every
// exercise is deleted and re-inserted from the submitted input. A field that
// insertTemplateBody does not copy is therefore not preserved across an edit, it
// is destroyed — which is why Create and Update share that one function.
func UpdateTemplate(id int64, name, focus string, circuits []TemplateCircuitInput, exercises []TemplateExerciseInput) (*Template, error) {
	if err := normalizeTemplateInput(name, circuits, exercises); err != nil {
		return nil, err
	}

	o := orm.NewOrm()
	t := &Template{ID: id}
	if err := o.Read(t); err != nil {
		return nil, errors.New("not found")
	}
	t.Name = name
	t.Focus = focus

	tx, err := o.Begin()
	if err != nil {
		return nil, err
	}

	if _, err := tx.Update(t); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Exercises first: they hold the foreign key into template_circuits.
	if _, err := tx.Raw("DELETE FROM template_exercises WHERE template_id = ?", id).Exec(); err != nil {
		tx.Rollback()
		return nil, err
	}

	if _, err := tx.Raw("DELETE FROM template_circuits WHERE template_id = ?", id).Exec(); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := insertTemplateBody(tx, id, circuits, exercises); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return t, nil
}

func GetTemplateByID(id int64) (*Template, []*TemplateExercise, error) {
	o := orm.NewOrm()
	t := &Template{ID: id}
	if err := o.Read(t); err != nil {
		return nil, nil, errors.New("not found")
	}
	var exercises []*TemplateExercise
	_, err := o.QueryTable(&TemplateExercise{}).Filter("TemplateID", id).OrderBy("SortOrder").All(&exercises)
	if err != nil {
		return nil, nil, err
	}
	return t, exercises, nil
}

// GetTemplateCircuits reads a template's circuits in display order. It is a
// separate read rather than a third return from GetTemplateByID because the
// session and exercise controllers call GetByID for the exercise list alone and
// have no use for circuits.
func GetTemplateCircuits(templateID int64) ([]*TemplateCircuit, error) {
	o := orm.NewOrm()
	var circuits []*TemplateCircuit
	_, err := o.QueryTable(&TemplateCircuit{}).Filter("TemplateID", templateID).OrderBy("SortOrder").All(&circuits)
	if err != nil {
		return nil, err
	}
	return circuits, nil
}
