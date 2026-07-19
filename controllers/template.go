package controllers

import (
	"fmt"
	"myGymPal/models"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

var blockOrder = []string{"main", "abs", "cardio", "stretch"}

var blockLabels = map[string]string{
	"main":    "Exercises",
	"abs":     "Abs",
	"cardio":  "Cardio",
	"stretch": "Stretch",
}

type templateExerciseBlock struct {
	Block     string
	Label     string
	Exercises []*models.TemplateExercise
}

// groupTemplateExercises buckets the loose exercises by block for the show page.
// Exercises inside a circuit are left out: the circuit renders as its own card,
// and listing them here as well would show every circuit exercise twice.
func groupTemplateExercises(exercises []*models.TemplateExercise) []templateExerciseBlock {
	byBlock := map[string][]*models.TemplateExercise{}
	for _, ex := range exercises {
		if ex.CircuitID != nil {
			continue
		}
		b := ex.Block
		if b == "" {
			b = "main"
		}
		byBlock[b] = append(byBlock[b], ex)
	}
	var blocks []templateExerciseBlock
	for _, key := range blockOrder {
		if exs, ok := byBlock[key]; ok {
			blocks = append(blocks, templateExerciseBlock{Block: key, Label: blockLabels[key], Exercises: exs})
		}
	}
	return blocks
}

// templateCircuitView is a circuit and its exercises, rendered on the show page
// as one card rather than as a run of loose exercises.
type templateCircuitView struct {
	Name              string
	Rounds            int
	TransitionSeconds int
	Exercises         []*models.TemplateExercise
}

func groupTemplateCircuits(circuits []*models.TemplateCircuit, exercises []*models.TemplateExercise) []templateCircuitView {
	byCircuit := map[int64][]*models.TemplateExercise{}
	for _, ex := range exercises {
		if ex.CircuitID != nil {
			byCircuit[*ex.CircuitID] = append(byCircuit[*ex.CircuitID], ex)
		}
	}
	views := make([]templateCircuitView, 0, len(circuits))
	for _, c := range circuits {
		views = append(views, templateCircuitView{
			Name:              c.Name,
			Rounds:            c.Rounds,
			TransitionSeconds: c.TransitionSeconds,
			Exercises:         byCircuit[c.ID],
		})
	}
	return views
}

type TemplateController struct {
	beego.Controller
}

// newFormChrome and editFormChrome supply the handful of values that differ
// between the create and edit renderings of partials/template_form.tpl. Each is
// called from both of its page's render paths — the initial GET and the
// validation-error re-render — because a c.Data key missed on one of them
// renders as an empty string rather than an error.
func (c *TemplateController) newFormChrome() {
	c.Data["BackURL"] = "/templates"
	c.Data["BackLabel"] = "Templates"
	c.Data["Heading"] = "New Workout Template"
	c.Data["FormAction"] = "/templates/new"
	c.Data["SubmitLabel"] = "Create Template"
}

func (c *TemplateController) editFormChrome(tmpl *models.Template) {
	showURL := fmt.Sprintf("/templates/%d", tmpl.ID)
	c.Data["BackURL"] = showURL
	c.Data["BackLabel"] = tmpl.Name
	c.Data["Heading"] = "Edit Template"
	c.Data["FormAction"] = showURL
	c.Data["SubmitLabel"] = "Save Changes"
}

// exerciseForm is a view-model for one exercise row on the form.
//
// Index is the number in the row's field names (exercise_name_3, block_3, and
// so on). It is carried explicitly rather than taken from the template's range
// index because circuit rows are rendered nested inside their circuit's card,
// where a range index would restart at zero and collide with the loose rows.
type exerciseForm struct {
	Index        int
	Name         string
	IsBodyweight bool
	IsTimeBased  bool
	Block        string
	WorkSeconds  int
}

// circuitForm is a view-model for one circuit card on the form, holding the
// exercise rows that belong to it.
type circuitForm struct {
	Index             int
	Name              string
	Rounds            int
	TransitionSeconds int
	Exercises         []exerciseForm
}

// templateFormBody is everything partials/template_form.tpl renders that is not
// page chrome. It exists so that the four render paths into that form — New,
// Create's error re-render, Edit and Update's error re-render — cannot disagree
// about which keys the template needs: they each build this struct and hand it
// to setFormBody, and a field left out is a compile error rather than a key that
// renders blank.
type templateFormBody struct {
	Name     string
	Focus    string
	Circuits []circuitForm
	// Loose holds the exercises that are not in any circuit.
	Loose []exerciseForm
	// ExerciseCount counts every row on the form, circuit rows included, because
	// it is what the server loops over when the form comes back.
	ExerciseCount int
}

func (c *TemplateController) setFormBody(b templateFormBody) {
	c.Data["Name"] = b.Name
	c.Data["Focus"] = b.Focus
	c.Data["Circuits"] = b.Circuits
	c.Data["Exercises"] = b.Loose
	c.Data["CircuitCount"] = len(b.Circuits)
	c.Data["ExerciseCount"] = b.ExerciseCount
}

// formInt reads an integer form field, treating anything unparseable as zero.
// The caller clamps it; the form is never trusted.
func (c *TemplateController) formInt(key string) int {
	n, err := strconv.Atoi(c.GetString(key))
	if err != nil {
		return 0
	}
	return n
}

// parseTemplateForm reads the circuit and exercise rows out of a submitted form.
// It returns the repository inputs alongside the view-model to re-render with if
// validation then fails.
//
// Both the circuits and the exercises are renumbered contiguously as they are
// read, because rows the user emptied out are dropped. That matters for circuits
// specifically: the exercises reference their circuit by position, so dropping
// circuit 1 would leave every exercise that named circuit 2 pointing at the
// wrong one. The submitted positions are mapped to the compacted ones, and an
// exercise whose circuit disappeared becomes a loose exercise rather than a
// dangling reference.
func (c *TemplateController) parseTemplateForm() ([]models.TemplateCircuitInput, []models.TemplateExerciseInput, templateFormBody) {
	name := c.GetString("name")
	focus := c.GetString("focus")

	circuitCount := c.formInt("circuit_count")
	circuits := make([]models.TemplateCircuitInput, 0, circuitCount)
	circuitForms := make([]circuitForm, 0, circuitCount)
	// submitted circuit position -> position after empty circuits are dropped
	circuitIndexMap := make(map[int]int, circuitCount)

	for i := 0; i < circuitCount; i++ {
		cName := c.GetString(fmt.Sprintf("circuit_name_%d", i))
		if cName == "" {
			continue
		}
		rounds := models.ValidRounds(c.formInt(fmt.Sprintf("circuit_rounds_%d", i)))
		transition := models.ValidSeconds(c.formInt(fmt.Sprintf("circuit_transition_%d", i)))

		idx := len(circuits)
		circuitIndexMap[i] = idx
		circuits = append(circuits, models.TemplateCircuitInput{
			Name:              cName,
			Rounds:            rounds,
			TransitionSeconds: transition,
			SortOrder:         idx,
		})
		circuitForms = append(circuitForms, circuitForm{
			Index:             idx,
			Name:              cName,
			Rounds:            rounds,
			TransitionSeconds: transition,
		})
	}

	exerciseCount := c.formInt("exercise_count")
	inputs := make([]models.TemplateExerciseInput, 0, exerciseCount)
	loose := make([]exerciseForm, 0, exerciseCount)

	for i := 0; i < exerciseCount; i++ {
		exName := c.GetString(fmt.Sprintf("exercise_name_%d", i))
		if exName == "" {
			continue
		}
		isBodyweight := c.GetString(fmt.Sprintf("is_bodyweight_%d", i)) != ""
		isTimeBased := c.GetString(fmt.Sprintf("is_time_based_%d", i)) != ""
		block := models.ValidBlock(c.GetString(fmt.Sprintf("block_%d", i)))
		workSeconds := models.ValidSeconds(c.formInt(fmt.Sprintf("work_seconds_%d", i)))

		// A row with no circuit_index, or one naming a circuit that was emptied
		// out of the form, is a loose exercise. Falling back to zero here instead
		// would silently file it under the first circuit.
		circuitIndex := models.NoCircuit
		if submitted, err := strconv.Atoi(c.GetString(fmt.Sprintf("circuit_index_%d", i))); err == nil {
			if mapped, ok := circuitIndexMap[submitted]; ok {
				circuitIndex = mapped
			}
		}

		idx := len(inputs)
		inputs = append(inputs, models.TemplateExerciseInput{
			Name:         exName,
			IsBodyweight: isBodyweight,
			IsTimeBased:  isTimeBased,
			Block:        block,
			SortOrder:    idx,
			CircuitIndex: circuitIndex,
			WorkSeconds:  workSeconds,
		})

		row := exerciseForm{
			Index:        idx,
			Name:         exName,
			IsBodyweight: isBodyweight,
			IsTimeBased:  isTimeBased,
			Block:        block,
			WorkSeconds:  workSeconds,
		}
		if circuitIndex == models.NoCircuit {
			loose = append(loose, row)
		} else {
			circuitForms[circuitIndex].Exercises = append(circuitForms[circuitIndex].Exercises, row)
		}
	}

	body := templateFormBody{
		Name:          name,
		Focus:         focus,
		Circuits:      circuitForms,
		Loose:         loose,
		ExerciseCount: len(inputs),
	}
	return circuits, inputs, body
}

// savedFormBody rebuilds the form view-model from what is in the database, for
// the edit page.
func savedFormBody(tmpl *models.Template, circuits []*models.TemplateCircuit, exercises []*models.TemplateExercise) templateFormBody {
	circuitForms := make([]circuitForm, len(circuits))
	positionByID := make(map[int64]int, len(circuits))
	for i, ci := range circuits {
		positionByID[ci.ID] = i
		circuitForms[i] = circuitForm{
			Index:             i,
			Name:              ci.Name,
			Rounds:            ci.Rounds,
			TransitionSeconds: ci.TransitionSeconds,
		}
	}

	loose := make([]exerciseForm, 0, len(exercises))
	for i, ex := range exercises {
		row := exerciseForm{
			Index:        i,
			Name:         ex.Name,
			IsBodyweight: ex.IsBodyweight,
			IsTimeBased:  ex.IsTimeBased,
			Block:        ex.Block,
			WorkSeconds:  ex.WorkSeconds,
		}
		if ex.CircuitID == nil {
			loose = append(loose, row)
			continue
		}
		// An exercise whose circuit is missing would otherwise vanish from the
		// form and be deleted by the next save. Show it as a loose exercise.
		pos, ok := positionByID[*ex.CircuitID]
		if !ok {
			loose = append(loose, row)
			continue
		}
		circuitForms[pos].Exercises = append(circuitForms[pos].Exercises, row)
	}

	return templateFormBody{
		Name:          tmpl.Name,
		Focus:         tmpl.Focus,
		Circuits:      circuitForms,
		Loose:         loose,
		ExerciseCount: len(exercises),
	}
}

func (c *TemplateController) Index() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	tmpls, err := Templates.GetAll()
	if err != nil {
		c.Redirect("/error", 302)
		return
	}

	flash := beego.ReadFromRequest(&c.Controller)
	if msg, ok := flash.Data["success"]; ok {
		c.Data["Success"] = msg
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "templates"
	c.Data["Templates"] = tmpls
	c.TplName = "templates/index.tpl"
}

func (c *TemplateController) New() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "templates"
	c.Data["ExerciseLibraryJSON"] = exerciseLibraryJSON(userID.(int64))
	c.setFormBody(templateFormBody{
		Loose:         []exerciseForm{{Index: 0, Block: "main"}},
		ExerciseCount: 1,
	})
	c.newFormChrome()
	c.TplName = "templates/new.tpl"
}

func (c *TemplateController) Create() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	circuits, inputs, body := c.parseTemplateForm()

	renderForm := func(errMsg string) {
		c.Data["LoggedIn"] = true
		c.Data["ActivePage"] = "templates"
		c.Data["Error"] = errMsg
		c.Data["ExerciseLibraryJSON"] = exerciseLibraryJSON(userID.(int64))
		c.setFormBody(body)
		c.newFormChrome()
		c.TplName = "templates/new.tpl"
	}

	if body.Name == "" {
		renderForm("Template name is required.")
		return
	}

	tmpl, err := Templates.Create(body.Name, body.Focus, circuits, inputs)
	if err != nil {
		renderForm(err.Error())
		return
	}

	flash := beego.NewFlash()
	flash.Success("Template created.")
	flash.Store(&c.Controller)
	c.Redirect(fmt.Sprintf("/templates/%d", tmpl.ID), 302)
}

func (c *TemplateController) Show() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Redirect("/templates", 302)
		return
	}

	tmpl, exercises, err := Templates.GetByID(id)
	if err != nil {
		c.Redirect("/templates", 302)
		return
	}

	circuits, err := Templates.GetCircuits(id)
	if err != nil {
		c.Redirect("/templates", 302)
		return
	}

	flash := beego.ReadFromRequest(&c.Controller)
	if msg, ok := flash.Data["success"]; ok {
		c.Data["Success"] = msg
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "templates"
	c.Data["Template"] = tmpl
	c.Data["Circuits"] = groupTemplateCircuits(circuits, exercises)
	c.Data["ExerciseBlocks"] = groupTemplateExercises(exercises)
	c.TplName = "templates/show.tpl"
}

func (c *TemplateController) Edit() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Redirect("/templates", 302)
		return
	}

	tmpl, exercises, err := Templates.GetByID(id)
	if err != nil {
		c.Redirect("/templates", 302)
		return
	}

	circuits, err := Templates.GetCircuits(id)
	if err != nil {
		c.Redirect("/templates", 302)
		return
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "templates"
	c.Data["Template"] = tmpl
	c.Data["ExerciseLibraryJSON"] = exerciseLibraryJSON(userID.(int64))
	c.setFormBody(savedFormBody(tmpl, circuits, exercises))
	c.editFormChrome(tmpl)
	c.TplName = "templates/edit.tpl"
}

func (c *TemplateController) Update() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Redirect("/templates", 302)
		return
	}

	tmpl, _, err := Templates.GetByID(id)
	if err != nil {
		c.Redirect("/templates", 302)
		return
	}

	circuits, inputs, body := c.parseTemplateForm()

	renderForm := func(errMsg string) {
		c.Data["LoggedIn"] = true
		c.Data["ActivePage"] = "templates"
		c.Data["Error"] = errMsg
		c.Data["Template"] = tmpl
		c.Data["ExerciseLibraryJSON"] = exerciseLibraryJSON(userID.(int64))
		c.setFormBody(body)
		c.editFormChrome(tmpl)
		c.TplName = "templates/edit.tpl"
	}

	if body.Name == "" {
		renderForm("Template name is required.")
		return
	}

	if _, err := Templates.Update(id, body.Name, body.Focus, circuits, inputs); err != nil {
		renderForm(err.Error())
		return
	}

	flash := beego.NewFlash()
	flash.Success("Template updated.")
	flash.Store(&c.Controller)
	c.Redirect(fmt.Sprintf("/templates/%d", id), 302)
}

func (c *TemplateController) Delete() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Redirect("/templates", 302)
		return
	}

	if err := Templates.Delete(id); err != nil {
		c.Redirect("/templates", 302)
		return
	}

	flash := beego.NewFlash()
	flash.Success("Template deleted.")
	flash.Store(&c.Controller)
	c.Redirect("/templates", 302)
}
