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

func (c *SessionController) New() {
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

	count, err := Sessions.CountByProgram(programID)
	if err != nil {
		c.Redirect("/error", 302)
		return
	}

	phase, week, workoutNum, _ := models.CalculateNextSession(count, program.WeeksPerPhase, program.WorkoutsPerWeek)

	templates, _ := Templates.GetAll()

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "programs"
	c.Data["Program"] = program
	c.Data["PhaseNumber"] = phase
	c.Data["WeekNumber"] = week
	c.Data["WorkoutNumber"] = workoutNum
	c.Data["Templates"] = templates
	c.TplName = "sessions/new.tpl"
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

	phaseNumber, err := strconv.Atoi(c.GetString("phase_number"))
	if err != nil || phaseNumber < 1 {
		c.Redirect(fmt.Sprintf("/programs/%d/sessions/new", programID), 302)
		return
	}

	weekNumber, err := strconv.Atoi(c.GetString("week_number"))
	if err != nil || weekNumber < 1 {
		c.Redirect(fmt.Sprintf("/programs/%d/sessions/new", programID), 302)
		return
	}

	workoutNumber, err := strconv.Atoi(c.GetString("workout_number"))
	if err != nil || workoutNumber < 1 {
		c.Redirect(fmt.Sprintf("/programs/%d/sessions/new", programID), 302)
		return
	}

	isDeload := weekNumber == program.WeeksPerPhase
	now := time.Now().UTC()

	session, err := Sessions.Create(programID, userID.(int64), phaseNumber, weekNumber, workoutNumber, isDeload, now.Truncate(24*time.Hour))
	if err != nil {
		c.Redirect("/error", 302)
		return
	}

	// Copy template exercises into the new session if a template was selected.
	if templateID, err := strconv.ParseInt(c.GetString("template_id"), 10, 64); err == nil && templateID > 0 {
		if _, exercises, err := Templates.GetByID(templateID); err == nil {
			for _, ex := range exercises {
				SessionExercises.Create(session.ID, ex.Name, ex.IsBodyweight, ex.GoalWeight, ex.WeightUnit)
			}
		}
	}

	c.Redirect(fmt.Sprintf("/sessions/%d", session.ID), 302)
}

func (c *SessionController) Delete() {
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

	// Need the program ID to redirect back; look up the session first.
	session, err := Sessions.GetByID(id, userID.(int64))
	if err != nil {
		c.Redirect("/programs", 302)
		return
	}
	programID := session.ProgramID

	if err := Sessions.Delete(id, userID.(int64)); err != nil {
		c.Redirect(fmt.Sprintf("/programs/%d", programID), 302)
		return
	}

	flash := beego.NewFlash()
	flash.Success("Workout deleted.")
	flash.Store(&c.Controller)
	c.Redirect(fmt.Sprintf("/programs/%d", programID), 302)
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

	exercises, err := SessionExercises.GetBySession(id)
	if err != nil {
		exercises = []*models.SessionExerciseView{}
	}

	user, err := Users.GetByID(userID.(int64))
	weightUnit := "lb"
	if err == nil {
		weightUnit = user.WeightUnit
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "programs"
	c.Data["Session"] = session
	c.Data["Program"] = program
	c.Data["Exercises"] = exercises
	c.Data["WeightUnit"] = weightUnit
	c.TplName = "sessions/show.tpl"
}

func (c *SessionController) AddExercise() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	sessionID, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Redirect("/programs", 302)
		return
	}

	// Verify session ownership.
	_, err = Sessions.GetByID(sessionID, userID.(int64))
	if err != nil {
		c.Redirect("/programs", 302)
		return
	}

	name := c.GetString("name")
	if name == "" {
		c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
		return
	}

	weightUnit := c.GetString("weight_unit")
	if weightUnit != "kg" {
		weightUnit = "lb"
	}

	goalWeightStr := c.GetString("goal_weight")
	goalWeight, _ := strconv.ParseFloat(goalWeightStr, 64)

	_, err = SessionExercises.Create(sessionID, name, false, goalWeight, weightUnit)
	if err != nil {
		c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
		return
	}

	c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
}

func (c *SessionController) LogSet() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	sessionID, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Redirect("/programs", 302)
		return
	}

	exerciseID, err := strconv.ParseInt(c.Ctx.Input.Param(":eid"), 10, 64)
	if err != nil {
		c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
		return
	}

	// Verify session ownership.
	_, err = Sessions.GetByID(sessionID, userID.(int64))
	if err != nil {
		c.Redirect("/programs", 302)
		return
	}

	// Verify the exercise belongs to this session.
	exercise, err := SessionExercises.GetByID(exerciseID)
	if err != nil || exercise.SessionID != sessionID {
		c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
		return
	}

	actualRepsStr := c.GetString("actual_reps")
	actualReps, err := strconv.Atoi(actualRepsStr)
	if err != nil || actualReps < 1 {
		c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
		return
	}

	weightUnit := c.GetString("weight_unit")
	if weightUnit != "kg" {
		weightUnit = "lb"
	}

	actualWeightStr := c.GetString("actual_weight")
	actualWeight, _ := strconv.ParseFloat(actualWeightStr, 64)

	count, err := SessionExercises.CountSetsByExercise(exerciseID)
	if err != nil {
		c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
		return
	}

	_, err = SessionExercises.LogSet(exerciseID, count+1, actualWeight, weightUnit, actualReps)
	if err != nil {
		c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
		return
	}

	c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
}
