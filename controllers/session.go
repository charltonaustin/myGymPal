package controllers

import (
	"fmt"
	"myGymPal/models"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type sessionExerciseBlock struct {
	Block     string
	Label     string
	Exercises []*models.SessionExerciseView
}

func groupSessionExercises(exercises []*models.SessionExerciseView) []sessionExerciseBlock {
	byBlock := map[string][]*models.SessionExerciseView{}
	for _, ev := range exercises {
		b := ev.Exercise.Block
		if b == "" {
			b = "main"
		}
		byBlock[b] = append(byBlock[b], ev)
	}
	var blocks []sessionExerciseBlock
	for _, key := range blockOrder {
		exs := byBlock[key]
		computeSupersetRuns(exs)
		// Always include the cardio block so the section is always visible.
		if exs != nil || key == "cardio" {
			blocks = append(blocks, sessionExerciseBlock{Block: key, Label: blockLabels[key], Exercises: exs})
		}
	}
	return blocks
}

// supersetMaxRun is the largest number of exercises one superset may hold.
const supersetMaxRun = 4

// computeSupersetRuns fills in SupersetLinked and SupersetLabel for one block's
// sort-ordered exercises.
//
// The raw LinkedToNext column is never trusted directly: an exercise is only
// *effectively* linked when there is a next exercise to flow into and the run it
// would extend is still under the cap. That is what makes a stale link on the last
// card of a block harmless — it is ignored, and the rest timer fires as normal.
//
// Runs of two or more get a per-block letter and a 1-based index (A1, A2, B1 …);
// a solo exercise gets no label.
func computeSupersetRuns(exercises []*models.SessionExerciseView) {
	runStart := 0
	for i, ev := range exercises {
		runLen := i - runStart + 1
		ev.SupersetLinked = ev.Exercise.LinkedToNext && i+1 < len(exercises) && runLen < supersetMaxRun
		ev.SupersetLabel = ""
		if !ev.SupersetLinked {
			runStart = i + 1
		}
	}

	letter := 'A'
	for start := 0; start < len(exercises); {
		end := start
		for end < len(exercises) && exercises[end].SupersetLinked {
			end++
		}
		// The run is exercises[start..end]: every member up to end links onward, and
		// end itself does not, so it is the last member.
		if end-start+1 >= 2 {
			for i := start; i <= end; i++ {
				exercises[i].SupersetLabel = fmt.Sprintf("%c%d", letter, i-start+1)
			}
			letter++
		}
		start = end + 1
	}
}

// supersetRunSize reports how many exercises would sit in the run containing the
// exercise at index i, if that exercise's link were turned on. It reads the raw
// links so that a request to link a 5th exercise into a full run can be rejected
// rather than silently capped.
func supersetRunSize(block []*models.SessionExerciseView, i int) int {
	start := i
	for start > 0 && block[start-1].Exercise.LinkedToNext {
		start--
	}
	end := i + 1
	for end < len(block)-1 && block[end].Exercise.LinkedToNext {
		end++
	}
	return end - start + 1
}

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

	var phase, week, workoutNum int
	logMode := c.GetString("sequential") == "1"
	if logMode {
		count, err := Sessions.CountByProgram(programID)
		if err != nil {
			logs.Error("SessionController.New: CountByProgram: %v", err)
			c.Redirect("/error", 302)
			return
		}
		phase, week, workoutNum, _ = models.CalculateNextSession(count, program.WeeksPerPhase, program.WorkoutsPerWeek)
	} else {
		last, err := Sessions.LatestByProgram(programID)
		if err != nil {
			logs.Error("SessionController.New: LatestByProgram: %v", err)
			c.Redirect("/error", 302)
			return
		}
		if last != nil {
			phase, week, workoutNum, _ = models.CalculateNextSessionFromLast(last, program.WeeksPerPhase, program.WorkoutsPerWeek)
		} else {
			phase, week, workoutNum = 1, 1, 1
		}
	}

	templates, _ := Templates.GetAll()

	var defaultTemplateID int64
	if wts, err := WorkoutTemplates.GetByProgram(programID); err == nil {
		for _, wt := range wts {
			if wt.WorkoutNumber == workoutNum {
				defaultTemplateID = wt.TemplateID
				break
			}
		}
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "programs"
	c.Data["Program"] = program
	c.Data["PhaseNumber"] = phase
	c.Data["WeekNumber"] = week
	c.Data["WorkoutNumber"] = workoutNum
	c.Data["Templates"] = templates
	c.Data["DefaultTemplateID"] = defaultTemplateID
	c.Data["DefaultDate"] = time.Now().Format("2006-01-02")
	c.Data["LogMode"] = logMode
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

	dateStr := c.GetString("date")
	sessionDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		sessionDate = time.Now().UTC().Truncate(24 * time.Hour)
	} else {
		sessionDate = sessionDate.UTC()
	}

	session, err := Sessions.Create(programID, userID.(int64), phaseNumber, weekNumber, workoutNumber, isDeload, sessionDate)
	if err != nil {
		logs.Error("SessionController.Create: %v", err)
		c.Redirect("/error", 302)
		return
	}

	// Copy template exercises into the new session if a template was selected.
	if templateID, err := strconv.ParseInt(c.GetString("template_id"), 10, 64); err == nil && templateID > 0 {
		if _, exercises, err := Templates.GetByID(templateID); err == nil {
			// Determine the rep minimum for the session's phase so it can be
			// stored as the goal reps on each exercise.
			goalReps := 0
			if phases, err := Phases.GetByProgram(programID); err == nil {
				for _, ph := range phases {
					if ph.PhaseNumber == phaseNumber {
						goalReps = ph.RepMin
						break
					}
				}
			}
			defaultUnit := "lb"
			if user, err := Users.GetByID(userID.(int64)); err == nil {
				defaultUnit = user.WeightUnit
			}
			for _, ex := range exercises {
				goalWeight := 0.0
				weightUnit := defaultUnit
				isTimeBased := false
				goalSeconds := 0
				exGoalReps := goalReps
				if libEx, err := Exercises.GetByName(userID.(int64), ex.Name); err == nil {
					goalWeight = libEx.GoalWeight
					weightUnit = libEx.WeightUnit
					isTimeBased = libEx.IsTimeBased
					goalSeconds = libEx.GoalSeconds
					if libEx.GoalRepMin > 0 {
						exGoalReps = libEx.GoalRepMin
					}
				}
				SessionExercises.Create(session.ID, ex.Name, ex.IsBodyweight, goalWeight, weightUnit, exGoalReps, ex.Block, isTimeBased, goalSeconds)
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

	// Enrich goal weight/unit/time/reps from the exercise library so the session always
	// reflects the current goal even if it was created before the library existed.
	for _, ev := range exercises {
		if libEx, err := Exercises.GetByName(userID.(int64), ev.Exercise.Name); err == nil {
			ev.Exercise.GoalWeight = libEx.GoalWeight
			ev.Exercise.WeightUnit = libEx.WeightUnit
			ev.Exercise.IsTimeBased = libEx.IsTimeBased
			ev.Exercise.GoalSeconds = libEx.GoalSeconds
			ev.Exercise.IsBodyweight = libEx.IsBodyweight
			if libEx.IsBodyweight {
				ev.GoalRepMin = libEx.GoalRepMin
				ev.GoalRepMax = libEx.GoalRepMax
			}
		}
	}

	user, err := Users.GetByID(userID.(int64))
	weightUnit := "lb"
	if err == nil {
		weightUnit = user.WeightUnit
	}

	// Convert logged set weights to each exercise's preferred unit.
	// Goal weight is already stored in the exercise's own WeightUnit — no conversion needed there.
	for _, ev := range exercises {
		if !ev.Exercise.IsBodyweight && !ev.Exercise.IsTimeBased {
			for _, s := range ev.Sets {
				s.ActualWeight = models.ConvertWeight(s.ActualWeight, s.WeightUnit, ev.Exercise.WeightUnit)
				s.WeightUnit = ev.Exercise.WeightUnit
			}
		}
	}

	// Find the rep range, default sets, and rest period for the session's phase.
	phaseRepMin, phaseRepMax, phaseDefaultSets, phaseRestSeconds := 0, 0, 0, 0
	if phases, err := Phases.GetByProgram(session.ProgramID); err == nil {
		for _, ph := range phases {
			if ph.PhaseNumber == session.PhaseNumber {
				phaseRepMin = ph.RepMin
				phaseRepMax = ph.RepMax
				phaseDefaultSets = ph.DefaultSets
				phaseRestSeconds = ph.RestSeconds
				break
			}
		}
	}

	// Mark HitMax on exercises where the user hit max reps at goal weight for all
	// required sets in the previous session for this program.
	if phaseRepMax > 0 {
		allSessions, err := Sessions.GetByProgram(session.ProgramID)
		if err == nil {
			var prevSessionID int64
			for i, s := range allSessions {
				if s.ID == id && i+1 < len(allSessions) {
					prevSessionID = allSessions[i+1].ID
					break
				}
			}
			if prevSessionID > 0 {
				if prevExs, err := SessionExercises.GetBySession(prevSessionID); err == nil {
					prevSetsByName := make(map[string][]*models.SessionSet, len(prevExs))
					for _, pev := range prevExs {
						prevSetsByName[pev.Exercise.Name] = pev.Sets
					}
					reqSets := phaseDefaultSets
					if reqSets < 1 {
						reqSets = 1
					}
					for _, ev := range exercises {
						if ev.Exercise.IsTimeBased {
							continue
						}
						prevSets := prevSetsByName[ev.Exercise.Name]
						if len(prevSets) < reqSets {
							continue
						}
						// For bodyweight exercises with a custom rep goal, use that;
						// otherwise fall back to the phase rep max.
						repMax := phaseRepMax
						if ev.Exercise.IsBodyweight && ev.GoalRepMax > 0 {
							repMax = ev.GoalRepMax
						}
						if repMax == 0 {
							continue
						}
						hitMax := true
						for _, s := range prevSets {
							if s.ActualReps < repMax {
								hitMax = false
								break
							}
							if !ev.Exercise.IsBodyweight && ev.Exercise.GoalWeight > 0 {
								convertedActual := models.ConvertWeight(s.ActualWeight, s.WeightUnit, ev.Exercise.WeightUnit)
								if convertedActual < ev.Exercise.GoalWeight {
									hitMax = false
									break
								}
							}
						}
						ev.HitMax = hitMax
						if !hitMax && !ev.Exercise.IsBodyweight && ev.Exercise.GoalWeight > 0 {
							for _, s := range prevSets {
								convertedActual := models.ConvertWeight(s.ActualWeight, s.WeightUnit, ev.Exercise.WeightUnit)
								if convertedActual < ev.Exercise.GoalWeight {
									ev.BelowGoal = true
									break
								}
							}
						}
					}
				}
			}
		}
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "programs"
	c.Data["Session"] = session
	c.Data["Program"] = program
	c.Data["ExerciseBlocks"] = groupSessionExercises(exercises)
	c.Data["WeightUnit"] = weightUnit
	c.Data["ExWeightUnit"] = weightUnit
	c.Data["PhaseRepMin"] = phaseRepMin
	c.Data["PhaseRepMax"] = phaseRepMax
	c.Data["PhaseRestSeconds"] = phaseRestSeconds
	c.Data["ExerciseLibraryJSON"] = exerciseLibraryJSON(userID.(int64))
	c.TplName = "sessions/show.tpl"
}

// UpdateRest persists a new default rest period for the session's phase so the
// change made from the in-session rest-timer control sticks for future sessions.
// Responds with JSON for the AJAX caller on the session page.
func (c *SessionController) UpdateRest() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]interface{}{"ok": false, "error": "not logged in"}
		c.ServeJSON()
		return
	}

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"ok": false, "error": "invalid session"}
		c.ServeJSON()
		return
	}

	restSeconds, err := strconv.Atoi(c.GetString("rest_seconds"))
	if err != nil || restSeconds < 0 {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"ok": false, "error": "invalid rest_seconds"}
		c.ServeJSON()
		return
	}

	session, err := Sessions.GetByID(id, userID.(int64))
	if err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"ok": false, "error": "session not found"}
		c.ServeJSON()
		return
	}

	// Confirm the session's program belongs to the user before touching its phase.
	if _, err := Programs.GetByID(session.ProgramID, userID.(int64)); err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"ok": false, "error": "program not found"}
		c.ServeJSON()
		return
	}

	if err := Phases.UpdateRestSeconds(session.ProgramID, session.PhaseNumber, restSeconds); err != nil {
		logs.Error("SessionController.UpdateRest: UpdateRestSeconds: %v", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{"ok": false, "error": "could not save"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{"ok": true, "rest_seconds": restSeconds}
	c.ServeJSON()
}

// UpdateLink toggles whether an exercise flows straight into the one below it in
// its block — a superset — instead of triggering the rest timer. Responds with
// JSON for the AJAX caller on the session page, which updates the card in place so
// a running rest timer is never lost to a reload.
func (c *SessionController) UpdateLink() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]interface{}{"ok": false, "error": "not logged in"}
		c.ServeJSON()
		return
	}

	sessionID, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"ok": false, "error": "invalid session"}
		c.ServeJSON()
		return
	}

	eid, err := strconv.ParseInt(c.Ctx.Input.Param(":eid"), 10, 64)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"ok": false, "error": "invalid exercise"}
		c.ServeJSON()
		return
	}

	// Ownership lives here, not in the repository: the session must belong to the
	// caller, and the exercise must belong to that session.
	if _, err := Sessions.GetByID(sessionID, userID.(int64)); err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"ok": false, "error": "session not found"}
		c.ServeJSON()
		return
	}

	exercise, err := SessionExercises.GetByID(eid)
	if err != nil || exercise.SessionID != sessionID {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]interface{}{"ok": false, "error": "exercise not found"}
		c.ServeJSON()
		return
	}

	linked := c.GetString("linked") == "true"

	// Turning a link off is always safe. Turning one on is only valid when there is
	// a next exercise in the same block and the run stays within the cap.
	if linked {
		exercises, err := SessionExercises.GetBySession(sessionID)
		if err != nil {
			logs.Error("SessionController.UpdateLink: GetBySession: %v", err)
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = map[string]interface{}{"ok": false, "error": "could not load session"}
			c.ServeJSON()
			return
		}

		block := blockMates(exercises, exercise)
		i := indexOfExercise(block, eid)
		if i < 0 || i == len(block)-1 {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = map[string]interface{}{"ok": false, "error": "nothing to link to"}
			c.ServeJSON()
			return
		}
		if supersetRunSize(block, i) > supersetMaxRun {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = map[string]interface{}{"ok": false, "error": fmt.Sprintf("a superset holds at most %d exercises", supersetMaxRun)}
			c.ServeJSON()
			return
		}
	}

	if err := SessionExercises.UpdateLinkedToNext(eid, linked); err != nil {
		logs.Error("SessionController.UpdateLink: %v", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{"ok": false, "error": "could not save"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{"ok": true, "linked": linked}
	c.ServeJSON()
}

// blockMates returns the exercises sharing a block with the given one, in sort
// order. GetBySession already returns exercises sorted, so filtering preserves it.
func blockMates(exercises []*models.SessionExerciseView, exercise *models.SessionExercise) []*models.SessionExerciseView {
	target := exercise.Block
	if target == "" {
		target = "main"
	}
	var block []*models.SessionExerciseView
	for _, ev := range exercises {
		b := ev.Exercise.Block
		if b == "" {
			b = "main"
		}
		if b == target {
			block = append(block, ev)
		}
	}
	return block
}

// indexOfExercise returns the position of the exercise within the block, or -1.
func indexOfExercise(block []*models.SessionExerciseView, eid int64) int {
	for i, ev := range block {
		if ev.Exercise.ID == eid {
			return i
		}
	}
	return -1
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

	isBodyweight := c.GetString("is_bodyweight") != ""
	isTimeBased := c.GetString("is_time_based") != ""

	weightUnit := c.GetString("weight_unit")
	if weightUnit != "kg" {
		weightUnit = "lb"
	}

	goalWeightStr := c.GetString("goal_weight")
	goalWeight, _ := strconv.ParseFloat(goalWeightStr, 64)
	goalH, _ := strconv.Atoi(c.GetString("goal_h"))
	goalM, _ := strconv.Atoi(c.GetString("goal_m"))
	goalS, _ := strconv.Atoi(c.GetString("goal_s"))
	goalSeconds := goalH*3600 + goalM*60 + goalS

	block := validBlock(c.GetString("block"))
	goalReps := 0
	if libEx, err := Exercises.GetByName(userID.(int64), name); err == nil {
		if libEx.GoalRepMin > 0 {
			goalReps = libEx.GoalRepMin
		}
	} else {
		Exercises.Create(userID.(int64), name, isBodyweight, goalWeight, weightUnit, isTimeBased, goalSeconds, 0, 0, block)
	}
	_, err = SessionExercises.Create(sessionID, name, isBodyweight, goalWeight, weightUnit, goalReps, block, isTimeBased, goalSeconds)
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

	var actualReps int
	var actualWeight float64
	var weightUnit string
	var actualSeconds int
	var activityType string

	if exercise.IsTimeBased {
		ah, _ := strconv.Atoi(c.GetString("actual_h"))
		am, _ := strconv.Atoi(c.GetString("actual_m"))
		as, _ := strconv.Atoi(c.GetString("actual_s"))
		actualSeconds = ah*3600 + am*60 + as
		if actualSeconds < 1 {
			c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
			return
		}
		activityType = c.GetString("activity_type")
		weightUnit = "lb"
	} else {
		actualReps, err = strconv.Atoi(c.GetString("actual_reps"))
		if err != nil || actualReps < 1 {
			c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
			return
		}
		weightUnit = c.GetString("weight_unit")
		if weightUnit != "kg" {
			weightUnit = "lb"
		}
		actualWeight, _ = strconv.ParseFloat(c.GetString("actual_weight"), 64)
	}

	count, err := SessionExercises.CountSetsByExercise(exerciseID)
	if err != nil {
		logs.Error("SessionController.LogSet: CountSetsByExercise: %v", err)
		c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
		return
	}

	set, err := SessionExercises.LogSet(exerciseID, count+1, actualWeight, weightUnit, actualReps, actualSeconds, activityType)
	if err != nil {
		logs.Error("SessionController.LogSet: LogSet: %v", err)
		c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
		return
	}

	// After 3rd set at or above goal reps AND at or above goal weight, update the exercise library entry.
	newCount := count + 1
	if !exercise.IsTimeBased && newCount >= 3 && actualReps >= exercise.GoalReps {
		if ex, err := Exercises.GetByName(userID.(int64), exercise.Name); err == nil {
			convertedGoal := models.ConvertWeight(ex.GoalWeight, ex.WeightUnit, weightUnit)
			if ex.GoalWeight == 0 || actualWeight >= convertedGoal {
				Exercises.UpdateGoalWeight(ex.ID, userID.(int64), actualWeight, weightUnit)
			}
		}
	}

	// AJAX callers get JSON back so the client can render the new row with a delete button.
	if c.Ctx.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		c.Data["json"] = map[string]interface{}{
			"id":             set.ID,
			"set_number":     set.SetNumber,
			"is_time_based":  exercise.IsTimeBased,
			"actual_seconds": actualSeconds,
			"activity_type":  activityType,
		}
		c.ServeJSON()
		return
	}

	c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
}

func (c *SessionController) AddCardioActivity() {
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

	if _, err := Sessions.GetByID(sessionID, userID.(int64)); err != nil {
		c.Redirect("/programs", 302)
		return
	}

	name := c.GetString("name")
	if name == "" {
		name = "cardio"
	}
	cardioType := c.GetString("cardio_type")
	goalDuration, _ := strconv.Atoi(c.GetString("goal_duration"))
	actualDuration, _ := strconv.Atoi(c.GetString("actual_duration"))

	ex, err := SessionExercises.Create(sessionID, name, false, 0, "lb", 0, "cardio", false, 0)
	if err != nil {
		logs.Error("SessionController.AddCardioActivity: Create: %v", err)
		c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
		return
	}

	if _, err := SessionExercises.LogCardio(ex.ID, cardioType, goalDuration, actualDuration); err != nil {
		logs.Error("SessionController.AddCardioActivity: LogCardio: %v", err)
	}

	c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
}

func (c *SessionController) LogCardio() {
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
	if _, err := Sessions.GetByID(sessionID, userID.(int64)); err != nil {
		c.Redirect("/programs", 302)
		return
	}

	// Verify the exercise belongs to this session.
	exercise, err := SessionExercises.GetByID(exerciseID)
	if err != nil || exercise.SessionID != sessionID {
		c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
		return
	}

	cardioType := c.GetString("cardio_type")
	goalDuration, _ := strconv.Atoi(c.GetString("goal_duration"))
	actualDuration, _ := strconv.Atoi(c.GetString("actual_duration"))

	if _, err := SessionExercises.LogCardio(exerciseID, cardioType, goalDuration, actualDuration); err != nil {
		logs.Error("SessionController.LogCardio: %v", err)
	}

	c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
}

func (c *SessionController) DeleteCardioLog() {
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

	logID, err := strconv.ParseInt(c.Ctx.Input.Param(":lid"), 10, 64)
	if err != nil {
		c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
		return
	}

	// Verify session ownership.
	if _, err := Sessions.GetByID(sessionID, userID.(int64)); err != nil {
		c.Redirect("/programs", 302)
		return
	}

	SessionExercises.DeleteCardioLog(logID)
	c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
}

func (c *SessionController) DeleteExercise() {
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
	if _, err := Sessions.GetByID(sessionID, userID.(int64)); err != nil {
		c.Redirect("/programs", 302)
		return
	}

	// Verify the exercise belongs to this session.
	exercise, err := SessionExercises.GetByID(exerciseID)
	if err != nil || exercise.SessionID != sessionID {
		c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
		return
	}

	SessionExercises.DeleteExercise(exerciseID)
	c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
}

func (c *SessionController) DeleteSet() {
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

	setID, err := strconv.ParseInt(c.Ctx.Input.Param(":sid"), 10, 64)
	if err != nil {
		c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
		return
	}

	// Verify session ownership.
	if _, err := Sessions.GetByID(sessionID, userID.(int64)); err != nil {
		c.Redirect("/programs", 302)
		return
	}

	// Verify the exercise belongs to this session.
	exercise, err := SessionExercises.GetByID(exerciseID)
	if err != nil || exercise.SessionID != sessionID {
		c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
		return
	}

	SessionExercises.DeleteSet(setID)

	if c.Ctx.Input.Header("X-Requested-With") == "XMLHttpRequest" {
		c.Data["json"] = map[string]string{"status": "ok"}
		c.ServeJSON()
		return
	}
	c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
}

func (c *SessionController) ReorderExercises() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Data["json"] = map[string]string{"error": "unauthorized"}
		c.ServeJSON()
		return
	}

	sessionID, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "invalid session"}
		c.ServeJSON()
		return
	}

	if _, err := Sessions.GetByID(sessionID, userID.(int64)); err != nil {
		c.Data["json"] = map[string]string{"error": "not found"}
		c.ServeJSON()
		return
	}

	raw := strings.TrimSpace(c.GetString("ids"))
	if raw == "" {
		c.Data["json"] = map[string]string{"ok": "1"}
		c.ServeJSON()
		return
	}

	parts := strings.Split(raw, ",")
	ids := make([]int64, 0, len(parts))
	for _, p := range parts {
		id, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64)
		if err != nil {
			c.Data["json"] = map[string]string{"error": "invalid id"}
			c.ServeJSON()
			return
		}
		ids = append(ids, id)
	}

	if err := SessionExercises.UpdateSortOrders(sessionID, ids); err != nil {
		logs.Error("ReorderExercises: %v", err)
		c.Data["json"] = map[string]string{"error": "failed to save order"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]string{"ok": "1"}
	c.ServeJSON()
}

// UpdateExerciseUnit handles AJAX requests to change the display unit for one exercise.
// Updates the exercise library entry by name so the preference persists across all sessions.
func (c *SessionController) UpdateExerciseUnit() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Ctx.Output.SetStatus(401)
		return
	}
	sessionID, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		return
	}
	if _, err := Sessions.GetByID(sessionID, userID.(int64)); err != nil {
		c.Ctx.Output.SetStatus(404)
		return
	}
	eid, err := strconv.ParseInt(c.Ctx.Input.Param(":eid"), 10, 64)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		return
	}
	newUnit := c.GetString("weight_unit")
	if newUnit != "lb" && newUnit != "kg" {
		c.Ctx.Output.SetStatus(400)
		return
	}
	se, err := SessionExercises.GetByID(eid)
	if err != nil || se.SessionID != sessionID {
		c.Ctx.Output.SetStatus(404)
		return
	}
	libEx, err := Exercises.GetByName(userID.(int64), se.Name)
	if err != nil {
		c.Ctx.Output.SetStatus(404)
		return
	}
	converted := models.ConvertWeight(libEx.GoalWeight, libEx.WeightUnit, newUnit)
	if err := Exercises.UpdateGoalWeight(libEx.ID, userID.(int64), converted, newUnit); err != nil {
		logs.Error("UpdateExerciseUnit: %v", err)
		c.Ctx.Output.SetStatus(500)
		return
	}
	c.Data["json"] = map[string]interface{}{"ok": true}
	c.ServeJSON()
}

func (c *SessionController) ChangeExercise() {
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
	if _, err := Sessions.GetByID(sessionID, userID.(int64)); err != nil {
		c.Redirect("/programs", 302)
		return
	}
	eid, err := strconv.ParseInt(c.Ctx.Input.Param(":eid"), 10, 64)
	if err != nil {
		c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
		return
	}
	se, err := SessionExercises.GetByID(eid)
	if err != nil || se.SessionID != sessionID {
		c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
		return
	}
	newName := strings.TrimSpace(c.GetString("name"))
	if newName == "" {
		c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
		return
	}
	if err := SessionExercises.UpdateName(eid, newName); err != nil {
		logs.Error("ChangeExercise: %v", err)
	}
	c.Redirect(fmt.Sprintf("/sessions/%d", sessionID), 302)
}
