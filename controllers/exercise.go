package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type ExerciseController struct {
	beego.Controller
}

func (c *ExerciseController) Index() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	exercises, err := Exercises.GetAllByUser(userID.(int64))
	if err != nil {
		logs.Error("ExerciseController.Index: GetAllByUser: %v", err)
		c.Redirect("/error", 302)
		return
	}

	flash := beego.ReadFromRequest(&c.Controller)
	if msg, ok := flash.Data["success"]; ok {
		c.Data["Success"] = msg
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "exercises"
	c.Data["Exercises"] = exercises
	c.TplName = "exercises/index.tpl"
}

func (c *ExerciseController) New() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	weightUnit := "lb"
	if user, err := Users.GetByID(userID.(int64)); err == nil {
		weightUnit = user.WeightUnit
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "exercises"
	c.Data["WeightUnit"] = weightUnit
	c.Data["ExWeightUnit"] = weightUnit
	c.TplName = "exercises/new.tpl"
}

func (c *ExerciseController) Create() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	weightUnit := "lb"
	if user, err := Users.GetByID(userID.(int64)); err == nil {
		weightUnit = user.WeightUnit
	}

	name := c.GetString("name")
	isBodyweight := c.GetString("is_bodyweight") != ""
	isTimeBased := c.GetString("is_time_based") != ""
	goalWeightStr := c.GetString("goal_weight")
	goalWeight, _ := strconv.ParseFloat(goalWeightStr, 64)
	goalH, _ := strconv.Atoi(c.GetString("goal_h"))
	goalM, _ := strconv.Atoi(c.GetString("goal_m"))
	goalS, _ := strconv.Atoi(c.GetString("goal_s"))
	goalSeconds := goalH*3600 + goalM*60 + goalS
	goalRepMin, _ := strconv.Atoi(c.GetString("goal_rep_min"))
	goalRepMax, _ := strconv.Atoi(c.GetString("goal_rep_max"))
	exWeightUnit := c.GetString("weight_unit")
	if exWeightUnit != "kg" {
		exWeightUnit = "lb"
	}

	renderForm := func(errMsg string) {
		c.Data["LoggedIn"] = true
		c.Data["ActivePage"] = "exercises"
		c.Data["Error"] = errMsg
		c.Data["Name"] = name
		c.Data["IsBodyweight"] = isBodyweight
		c.Data["IsTimeBased"] = isTimeBased
		c.Data["GoalWeight"] = goalWeightStr
		c.Data["GoalHours"] = goalH
		c.Data["GoalMinutes"] = goalM
		c.Data["GoalSecsRemainder"] = goalS
		c.Data["GoalRepMin"] = goalRepMin
		c.Data["GoalRepMax"] = goalRepMax
		c.Data["WeightUnit"] = weightUnit
		c.Data["ExWeightUnit"] = exWeightUnit
		c.TplName = "exercises/new.tpl"
	}

	if name == "" {
		renderForm("Exercise name is required.")
		return
	}

	if _, err := Exercises.Create(userID.(int64), name, isBodyweight, goalWeight, exWeightUnit, isTimeBased, goalSeconds, goalRepMin, goalRepMax); err != nil {
		logs.Error("ExerciseController.Create: %v", err)
		renderForm(err.Error())
		return
	}

	flash := beego.NewFlash()
	flash.Success("%s added to your exercise library.", name)
	flash.Store(&c.Controller)
	c.Redirect("/exercises", 302)
}

func (c *ExerciseController) Edit() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Redirect("/exercises", 302)
		return
	}

	ex, err := Exercises.GetByID(id, userID.(int64))
	if err != nil {
		c.Redirect("/exercises", 302)
		return
	}

	weightUnit := "lb"
	if user, err := Users.GetByID(userID.(int64)); err == nil {
		weightUnit = user.WeightUnit
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "exercises"
	c.Data["WeightUnit"] = weightUnit
	c.Data["Exercise"] = ex
	c.Data["Name"] = ex.Name
	c.Data["IsBodyweight"] = ex.IsBodyweight
	c.Data["IsTimeBased"] = ex.IsTimeBased
	c.Data["GoalWeight"] = fmt.Sprintf("%g", ex.GoalWeight)
	c.Data["GoalHours"] = ex.GoalSeconds / 3600
	c.Data["GoalMinutes"] = (ex.GoalSeconds % 3600) / 60
	c.Data["GoalSecsRemainder"] = ex.GoalSeconds % 60
	c.Data["GoalRepMin"] = ex.GoalRepMin
	c.Data["GoalRepMax"] = ex.GoalRepMax
	c.Data["ExWeightUnit"] = ex.WeightUnit
	c.TplName = "exercises/edit.tpl"
}

func (c *ExerciseController) Update() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Redirect("/exercises", 302)
		return
	}

	ex, err := Exercises.GetByID(id, userID.(int64))
	if err != nil {
		c.Redirect("/exercises", 302)
		return
	}

	weightUnit := "lb"
	if user, err := Users.GetByID(userID.(int64)); err == nil {
		weightUnit = user.WeightUnit
	}

	name := c.GetString("name")
	isBodyweight := c.GetString("is_bodyweight") != ""
	isTimeBased := c.GetString("is_time_based") != ""
	goalWeightStr := c.GetString("goal_weight")
	goalWeight, _ := strconv.ParseFloat(goalWeightStr, 64)
	goalH, _ := strconv.Atoi(c.GetString("goal_h"))
	goalM, _ := strconv.Atoi(c.GetString("goal_m"))
	goalS, _ := strconv.Atoi(c.GetString("goal_s"))
	goalSeconds := goalH*3600 + goalM*60 + goalS
	goalRepMin, _ := strconv.Atoi(c.GetString("goal_rep_min"))
	goalRepMax, _ := strconv.Atoi(c.GetString("goal_rep_max"))
	exWeightUnit := c.GetString("weight_unit")
	if exWeightUnit != "kg" {
		exWeightUnit = "lb"
	}

	renderForm := func(errMsg string) {
		c.Data["LoggedIn"] = true
		c.Data["ActivePage"] = "exercises"
		c.Data["Error"] = errMsg
		c.Data["WeightUnit"] = weightUnit
		c.Data["Exercise"] = ex
		c.Data["Name"] = name
		c.Data["IsBodyweight"] = isBodyweight
		c.Data["IsTimeBased"] = isTimeBased
		c.Data["GoalWeight"] = goalWeightStr
		c.Data["GoalHours"] = goalH
		c.Data["GoalMinutes"] = goalM
		c.Data["GoalSecsRemainder"] = goalS
		c.Data["GoalRepMin"] = goalRepMin
		c.Data["GoalRepMax"] = goalRepMax
		c.Data["ExWeightUnit"] = exWeightUnit
		c.TplName = "exercises/edit.tpl"
	}

	if name == "" {
		renderForm("Exercise name is required.")
		return
	}

	if _, err := Exercises.Update(id, userID.(int64), name, isBodyweight, goalWeight, exWeightUnit, isTimeBased, goalSeconds, goalRepMin, goalRepMax); err != nil {
		logs.Error("ExerciseController.Update: %v", err)
		renderForm(err.Error())
		return
	}

	flash := beego.NewFlash()
	flash.Success("Exercise updated.")
	flash.Store(&c.Controller)
	c.Redirect("/exercises", 302)
}

func (c *ExerciseController) Delete() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Redirect("/exercises", 302)
		return
	}

	if err := Exercises.Delete(id, userID.(int64)); err != nil {
		c.Redirect("/exercises", 302)
		return
	}

	flash := beego.NewFlash()
	flash.Success("Exercise deleted.")
	flash.Store(&c.Controller)
	c.Redirect("/exercises", 302)
}

// UpdateGoalWeightJSON handles AJAX requests to update an exercise's goal weight.
// Looks up the exercise by name (enforcing ownership), then updates the goal weight and unit.
func (c *ExerciseController) UpdateGoalWeightJSON() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Data["json"] = map[string]string{"error": "unauthenticated"}
		c.ServeJSON()
		return
	}

	name := c.GetString("name")
	goalWeightStr := c.GetString("goal_weight")
	goalWeight, err := strconv.ParseFloat(goalWeightStr, 64)
	if err != nil || goalWeight < 0 {
		c.Data["json"] = map[string]string{"error": "invalid goal weight"}
		c.ServeJSON()
		return
	}
	weightUnit := c.GetString("weight_unit")
	if weightUnit != "kg" {
		weightUnit = "lb"
	}

	libEx, err := Exercises.GetByName(userID.(int64), name)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "exercise not found"}
		c.ServeJSON()
		return
	}

	if _, err := Exercises.Update(libEx.ID, userID.(int64), libEx.Name, libEx.IsBodyweight, goalWeight, weightUnit, libEx.IsTimeBased, libEx.GoalSeconds, libEx.GoalRepMin, libEx.GoalRepMax); err != nil {
		logs.Error("ExerciseController.UpdateGoalWeightJSON: %v", err)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{"ok": true, "goal_weight": goalWeight, "weight_unit": weightUnit}
	c.ServeJSON()
}

// exerciseLibraryJSON fetches the user's exercise library and returns a template.JS
// value safe for direct embedding in a <script> tag without HTML escaping.
func exerciseLibraryJSON(userID int64) template.JS {
	exercises, err := Exercises.GetAllByUser(userID)
	if err != nil || len(exercises) == 0 {
		return "[]"
	}
	type libEntry struct {
		Name         string  `json:"name"`
		GoalWeight   float64 `json:"goalWeight"`
		WeightUnit   string  `json:"weightUnit"`
		IsBodyweight bool    `json:"isBodyweight"`
		IsTimeBased  bool    `json:"isTimeBased"`
		GoalSeconds  int     `json:"goalSeconds"`
	}
	entries := make([]libEntry, len(exercises))
	for i, ex := range exercises {
		entries[i] = libEntry{
			Name:         ex.Name,
			GoalWeight:   ex.GoalWeight,
			WeightUnit:   ex.WeightUnit,
			IsBodyweight: ex.IsBodyweight,
			IsTimeBased:  ex.IsTimeBased,
			GoalSeconds:  ex.GoalSeconds,
		}
	}
	b, err := json.Marshal(entries)
	if err != nil {
		return "[]"
	}
	return template.JS(b)
}
