package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"

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
	goalWeightStr := c.GetString("goal_weight")
	goalWeight, _ := strconv.ParseFloat(goalWeightStr, 64)
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
		c.Data["GoalWeight"] = goalWeightStr
		c.Data["WeightUnit"] = weightUnit
		c.Data["ExWeightUnit"] = exWeightUnit
		c.TplName = "exercises/new.tpl"
	}

	if name == "" {
		renderForm("Exercise name is required.")
		return
	}

	if _, err := Exercises.Create(userID.(int64), name, isBodyweight, goalWeight, exWeightUnit); err != nil {
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
	c.Data["GoalWeight"] = fmt.Sprintf("%g", ex.GoalWeight)
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
	goalWeightStr := c.GetString("goal_weight")
	goalWeight, _ := strconv.ParseFloat(goalWeightStr, 64)
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
		c.Data["GoalWeight"] = goalWeightStr
		c.Data["ExWeightUnit"] = exWeightUnit
		c.TplName = "exercises/edit.tpl"
	}

	if name == "" {
		renderForm("Exercise name is required.")
		return
	}

	if _, err := Exercises.Update(id, userID.(int64), name, isBodyweight, goalWeight, exWeightUnit); err != nil {
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
	}
	entries := make([]libEntry, len(exercises))
	for i, ex := range exercises {
		entries[i] = libEntry{
			Name:         ex.Name,
			GoalWeight:   ex.GoalWeight,
			WeightUnit:   ex.WeightUnit,
			IsBodyweight: ex.IsBodyweight,
		}
	}
	b, err := json.Marshal(entries)
	if err != nil {
		return "[]"
	}
	return template.JS(b)
}
