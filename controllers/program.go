package controllers

import (
	"fmt"
	"myGymPal/models"
	"strconv"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

type ProgramController struct {
	beego.Controller
}

func (c *ProgramController) Index() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	programs, err := Programs.GetAllByUser(userID.(int64))
	if err != nil {
		c.Redirect("/error", 302)
		return
	}

	flash := beego.ReadFromRequest(&c.Controller)
	if msg, ok := flash.Data["success"]; ok {
		c.Data["Success"] = msg
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "programs"
	c.Data["Programs"] = programs
	c.TplName = "programs/index.tpl"
}

func (c *ProgramController) New() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "programs"
	c.Data["WeeksPerPhase"] = "8"
	c.Data["DefaultRepMin"] = ""
	c.Data["DefaultRepMax"] = ""
	c.TplName = "programs/new.tpl"
}

func (c *ProgramController) Create() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	name := c.GetString("name")
	startDateStr := c.GetString("start_date")
	numPhasesStr := c.GetString("num_phases")
	weeksPerPhaseStr := c.GetString("weeks_per_phase")
	defaultRepMinStr := c.GetString("default_rep_min")
	defaultRepMaxStr := c.GetString("default_rep_max")

	renderForm := func(errMsg string) {
		c.Data["LoggedIn"] = true
		c.Data["ActivePage"] = "programs"
		c.Data["Error"] = errMsg
		c.Data["Name"] = name
		c.Data["StartDate"] = startDateStr
		c.Data["NumPhases"] = numPhasesStr
		c.Data["WeeksPerPhase"] = weeksPerPhaseStr
		c.Data["DefaultRepMin"] = defaultRepMinStr
		c.Data["DefaultRepMax"] = defaultRepMaxStr
		c.TplName = "programs/new.tpl"
	}

	if name == "" {
		renderForm("Program name is required.")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		renderForm("Start date is required and must be a valid date (YYYY-MM-DD).")
		return
	}

	numPhases, err := strconv.Atoi(numPhasesStr)
	if err != nil || numPhases <= 0 {
		renderForm("Number of phases must be a positive number.")
		return
	}

	weeksPerPhase, err := strconv.Atoi(weeksPerPhaseStr)
	if err != nil || weeksPerPhase <= 0 {
		renderForm("Weeks per phase must be a positive number.")
		return
	}

	defaultRepMin, err := strconv.Atoi(defaultRepMinStr)
	if err != nil || defaultRepMin <= 0 {
		renderForm("Default min reps must be a positive number.")
		return
	}

	defaultRepMax, err := strconv.Atoi(defaultRepMaxStr)
	if err != nil || defaultRepMax < defaultRepMin {
		renderForm("Default max reps must be at least the min reps.")
		return
	}

	if _, err := Programs.Create(userID.(int64), name, startDate, numPhases, weeksPerPhase, defaultRepMin, defaultRepMax); err != nil {
		renderForm("Something went wrong. Please try again.")
		return
	}

	flash := beego.NewFlash()
	flash.Success("Program created.")
	flash.Store(&c.Controller)
	c.Redirect("/programs", 302)
}

func (c *ProgramController) Show() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Redirect("/programs", 302)
		return
	}

	program, err := Programs.GetByID(id, userID.(int64))
	if err != nil {
		c.Redirect("/programs", 302)
		return
	}

	phases, err := Phases.GetByProgram(id)
	if err != nil {
		c.Redirect("/error", 302)
		return
	}

	flash := beego.ReadFromRequest(&c.Controller)
	if msg, ok := flash.Data["success"]; ok {
		c.Data["Success"] = msg
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "programs"
	c.Data["Program"] = program
	c.Data["Phases"] = phases
	c.TplName = "programs/show.tpl"
}

func (c *ProgramController) UpdatePhases() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Redirect("/programs", 302)
		return
	}

	program, err := Programs.GetByID(id, userID.(int64))
	if err != nil {
		c.Redirect("/programs", 302)
		return
	}

	phases, err := Phases.GetByProgram(id)
	if err != nil {
		c.Redirect("/error", 302)
		return
	}

	renderShow := func(errMsg string, viewPhases []*models.Phase) {
		c.Data["LoggedIn"] = true
		c.Data["ActivePage"] = "programs"
		c.Data["Program"] = program
		c.Data["Phases"] = viewPhases
		c.Data["Error"] = errMsg
		c.TplName = "programs/show.tpl"
	}

	viewPhases := make([]*models.Phase, len(phases))
	updates := make([]models.PhaseUpdate, len(phases))
	for i, ph := range phases {
		repMinStr := c.GetString(fmt.Sprintf("rep_min_%d", ph.PhaseNumber))
		repMaxStr := c.GetString(fmt.Sprintf("rep_max_%d", ph.PhaseNumber))

		repMin, _ := strconv.Atoi(repMinStr)
		repMax, _ := strconv.Atoi(repMaxStr)

		viewPhases[i] = &models.Phase{
			ID:          ph.ID,
			ProgramID:   ph.ProgramID,
			PhaseNumber: ph.PhaseNumber,
			RepMin:      repMin,
			RepMax:      repMax,
		}
		updates[i] = models.PhaseUpdate{PhaseNumber: ph.PhaseNumber, RepMin: repMin, RepMax: repMax}
	}

	if err := Phases.UpdateRepRanges(id, updates); err != nil {
		renderShow(err.Error(), viewPhases)
		return
	}

	flash := beego.NewFlash()
	flash.Success("Rep ranges saved.")
	flash.Store(&c.Controller)
	c.Redirect(fmt.Sprintf("/programs/%d", id), 302)
}
