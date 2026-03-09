package controllers

import (
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

	programs, err := models.GetProgramsByUserID(userID.(int64))
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

	renderForm := func(errMsg string) {
		c.Data["LoggedIn"] = true
		c.Data["ActivePage"] = "programs"
		c.Data["Error"] = errMsg
		c.Data["Name"] = name
		c.Data["StartDate"] = startDateStr
		c.Data["NumPhases"] = numPhasesStr
		c.Data["WeeksPerPhase"] = weeksPerPhaseStr
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

	if _, err := models.CreateProgram(userID.(int64), name, startDate, numPhases, weeksPerPhase); err != nil {
		renderForm("Something went wrong. Please try again.")
		return
	}

	flash := beego.NewFlash()
	flash.Success("Program created.")
	flash.Store(&c.Controller)
	c.Redirect("/programs", 302)
}
