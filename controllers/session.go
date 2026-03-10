package controllers

import (
	"fmt"
	"myGymPal/models"
	"strconv"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

type SessionController struct {
	beego.Controller
}

func (c *SessionController) Create() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	programID, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Redirect("/programs", 302)
		return
	}

	program, err := Programs.GetByID(programID, userID.(int64))
	if err != nil {
		c.Redirect("/programs", 302)
		return
	}

	now := time.Now().UTC()
	phase, week, isDeload := models.CalculatePhaseAndWeek(now, program.StartDate, program.WeeksPerPhase)

	count, err := Sessions.CountByProgram(programID)
	if err != nil {
		c.Redirect("/error", 302)
		return
	}

	session, err := Sessions.Create(programID, userID.(int64), phase, week, count+1, isDeload, now.Truncate(24*time.Hour))
	if err != nil {
		c.Redirect("/error", 302)
		return
	}

	c.Redirect(fmt.Sprintf("/sessions/%d", session.ID), 302)
}

func (c *SessionController) Show() {
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

	session, err := Sessions.GetByID(id, userID.(int64))
	if err != nil {
		c.Redirect("/programs", 302)
		return
	}

	program, err := Programs.GetByID(session.ProgramID, userID.(int64))
	if err != nil {
		c.Redirect("/programs", 302)
		return
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "programs"
	c.Data["Session"] = session
	c.Data["Program"] = program
	c.TplName = "sessions/show.tpl"
}
