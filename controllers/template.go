package controllers

import (
	"fmt"
	"myGymPal/models"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

func exercisesToForms(exercises []*models.TemplateExercise) []exerciseForm {
	forms := make([]exerciseForm, len(exercises))
	for i, ex := range exercises {
		forms[i] = exerciseForm{
			Name:         ex.Name,
			IsBodyweight: ex.IsBodyweight,
			GoalWeight:   fmt.Sprintf("%g", ex.GoalWeight),
			WeightUnit:   ex.WeightUnit,
			RepMin:       strconv.Itoa(ex.RepMin),
			RepMax:       strconv.Itoa(ex.RepMax),
		}
	}
	return forms
}

type TemplateController struct {
	beego.Controller
}

// exerciseForm is a view-model for re-rendering exercise rows on validation error.
type exerciseForm struct {
	Name         string
	IsBodyweight bool
	GoalWeight   string
	WeightUnit   string
	RepMin       string
	RepMax       string
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

	user, err := Users.GetByID(userID.(int64))
	if err != nil {
		c.Redirect("/error", 302)
		return
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "templates"
	c.Data["WeightUnit"] = user.WeightUnit
	c.Data["Exercises"] = []exerciseForm{{WeightUnit: user.WeightUnit}}
	c.Data["ExerciseLibraryJSON"] = exerciseLibraryJSON(userID.(int64))
	c.TplName = "templates/new.tpl"
}

func (c *TemplateController) Create() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	user, userErr := Users.GetByID(userID.(int64))
	weightUnit := "lb"
	if userErr == nil {
		weightUnit = user.WeightUnit
	}

	name := c.GetString("name")
	focus := c.GetString("focus")
	countStr := c.GetString("exercise_count")
	count, _ := strconv.Atoi(countStr)

	forms := make([]exerciseForm, 0, count)
	inputs := make([]models.TemplateExerciseInput, 0, count)
	for i := 0; i < count; i++ {
		exName := c.GetString(fmt.Sprintf("exercise_name_%d", i))
		if exName == "" {
			continue // skip removed rows
		}
		isBodyweight := c.GetString(fmt.Sprintf("is_bodyweight_%d", i)) != ""
		goalWeightStr := c.GetString(fmt.Sprintf("goal_weight_%d", i))
		exWeightUnit := c.GetString(fmt.Sprintf("weight_unit_%d", i))
		if exWeightUnit == "" {
			exWeightUnit = weightUnit
		}
		repMinStr := c.GetString(fmt.Sprintf("rep_min_%d", i))
		repMaxStr := c.GetString(fmt.Sprintf("rep_max_%d", i))

		goalWeight, _ := strconv.ParseFloat(goalWeightStr, 64)
		repMin, _ := strconv.Atoi(repMinStr)
		repMax, _ := strconv.Atoi(repMaxStr)

		forms = append(forms, exerciseForm{
			Name:         exName,
			IsBodyweight: isBodyweight,
			GoalWeight:   goalWeightStr,
			WeightUnit:   exWeightUnit,
			RepMin:       repMinStr,
			RepMax:       repMaxStr,
		})
		inputs = append(inputs, models.TemplateExerciseInput{
			Name:         exName,
			IsBodyweight: isBodyweight,
			GoalWeight:   goalWeight,
			WeightUnit:   exWeightUnit,
			RepMin:       repMin,
			RepMax:       repMax,
			SortOrder:    len(inputs),
		})
	}

	renderForm := func(errMsg string) {
		c.Data["LoggedIn"] = true
		c.Data["ActivePage"] = "templates"
		c.Data["Error"] = errMsg
		c.Data["Name"] = name
		c.Data["Focus"] = focus
		c.Data["WeightUnit"] = weightUnit
		c.Data["Exercises"] = forms
		c.Data["ExerciseLibraryJSON"] = exerciseLibraryJSON(userID.(int64))
		c.TplName = "templates/new.tpl"
	}

	if name == "" {
		renderForm("Template name is required.")
		return
	}

	tmpl, err := Templates.Create(name, focus, inputs)
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

	flash := beego.ReadFromRequest(&c.Controller)
	if msg, ok := flash.Data["success"]; ok {
		c.Data["Success"] = msg
	}

	weightUnit := "lb"
	if user, err := Users.GetByID(userID.(int64)); err == nil {
		weightUnit = user.WeightUnit
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "templates"
	c.Data["WeightUnit"] = weightUnit
	c.Data["Template"] = tmpl
	c.Data["Exercises"] = exercises
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

	weightUnit := "lb"
	if user, err := Users.GetByID(userID.(int64)); err == nil {
		weightUnit = user.WeightUnit
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "templates"
	c.Data["WeightUnit"] = weightUnit
	c.Data["Template"] = tmpl
	c.Data["Name"] = tmpl.Name
	c.Data["Focus"] = tmpl.Focus
	c.Data["Exercises"] = exercisesToForms(exercises)
	c.Data["ExerciseLibraryJSON"] = exerciseLibraryJSON(userID.(int64))
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

	user, userErr := Users.GetByID(userID.(int64))
	weightUnit := "lb"
	if userErr == nil {
		weightUnit = user.WeightUnit
	}

	name := c.GetString("name")
	focus := c.GetString("focus")
	countStr := c.GetString("exercise_count")
	count, _ := strconv.Atoi(countStr)

	forms := make([]exerciseForm, 0, count)
	inputs := make([]models.TemplateExerciseInput, 0, count)
	for i := 0; i < count; i++ {
		exName := c.GetString(fmt.Sprintf("exercise_name_%d", i))
		if exName == "" {
			continue
		}
		isBodyweight := c.GetString(fmt.Sprintf("is_bodyweight_%d", i)) != ""
		goalWeightStr := c.GetString(fmt.Sprintf("goal_weight_%d", i))
		exWeightUnit := c.GetString(fmt.Sprintf("weight_unit_%d", i))
		if exWeightUnit == "" {
			exWeightUnit = weightUnit
		}
		repMinStr := c.GetString(fmt.Sprintf("rep_min_%d", i))
		repMaxStr := c.GetString(fmt.Sprintf("rep_max_%d", i))

		goalWeight, _ := strconv.ParseFloat(goalWeightStr, 64)
		repMin, _ := strconv.Atoi(repMinStr)
		repMax, _ := strconv.Atoi(repMaxStr)

		forms = append(forms, exerciseForm{
			Name:         exName,
			IsBodyweight: isBodyweight,
			GoalWeight:   goalWeightStr,
			WeightUnit:   exWeightUnit,
			RepMin:       repMinStr,
			RepMax:       repMaxStr,
		})
		inputs = append(inputs, models.TemplateExerciseInput{
			Name:         exName,
			IsBodyweight: isBodyweight,
			GoalWeight:   goalWeight,
			WeightUnit:   exWeightUnit,
			RepMin:       repMin,
			RepMax:       repMax,
			SortOrder:    len(inputs),
		})
	}

	renderForm := func(errMsg string) {
		c.Data["LoggedIn"] = true
		c.Data["ActivePage"] = "templates"
		c.Data["Error"] = errMsg
		c.Data["Template"] = tmpl
		c.Data["Name"] = name
		c.Data["Focus"] = focus
		c.Data["WeightUnit"] = weightUnit
		c.Data["Exercises"] = forms
		c.Data["ExerciseLibraryJSON"] = exerciseLibraryJSON(userID.(int64))
		c.TplName = "templates/edit.tpl"
	}

	if name == "" {
		renderForm("Template name is required.")
		return
	}

	if _, err := Templates.Update(id, name, focus, inputs); err != nil {
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
