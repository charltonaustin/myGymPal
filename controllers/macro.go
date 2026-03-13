package controllers

import (
	"encoding/json"
	"html/template"
	"myGymPal/models"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type MacroController struct {
	beego.Controller
}

type macroDay struct {
	Date     time.Time
	Entries  []*models.MacroEntry
	Protein  float64
	Carbs    float64
	Fat      float64
	Calories float64
}

type macroSummaryRow struct {
	Actual float64
	Goal   float64
	Pct    int
	AtGoal bool
}

type macroSummary struct {
	Days     int
	Protein  macroSummaryRow
	Carbs    macroSummaryRow
	Fat      macroSummaryRow
	Calories macroSummaryRow
	HasGoal  bool
}

func buildMacroSummary(days []macroDay, goal *models.MacroGoal) *macroSummary {
	n := len(days)
	if n == 0 {
		return nil
	}
	if n > 3 {
		n = 3
	}
	var p, c, f float64
	for i := 0; i < n; i++ {
		p += days[i].Protein
		c += days[i].Carbs
		f += days[i].Fat
	}
	dn := float64(n)
	s := &macroSummary{
		Days:     n,
		Protein:  macroSummaryRow{Actual: p / dn},
		Carbs:    macroSummaryRow{Actual: c / dn},
		Fat:      macroSummaryRow{Actual: f / dn},
		Calories: macroSummaryRow{Actual: (4*p + 4*c + 9*f) / dn},
	}
	if goal != nil {
		s.HasGoal = true
		pct := func(actual, goal float64) (int, bool) {
			if goal <= 0 {
				return 0, true
			}
			p := int(actual / goal * 100)
			return p, p >= 100
		}
		goalCal := 4*goal.Protein + 4*goal.Carbs + 9*goal.Fat
		s.Protein.Goal = goal.Protein
		s.Protein.Pct, s.Protein.AtGoal = pct(s.Protein.Actual, goal.Protein)
		s.Carbs.Goal = goal.Carbs
		s.Carbs.Pct, s.Carbs.AtGoal = pct(s.Carbs.Actual, goal.Carbs)
		s.Fat.Goal = goal.Fat
		s.Fat.Pct, s.Fat.AtGoal = pct(s.Fat.Actual, goal.Fat)
		s.Calories.Goal = goalCal
		s.Calories.Pct, s.Calories.AtGoal = pct(s.Calories.Actual, goalCal)
	}
	return s
}

func groupMacrosByDay(entries []*models.MacroEntry) []macroDay {
	var days []macroDay
	index := map[string]int{}
	for _, e := range entries {
		key := e.Date.Format("2006-01-02")
		if i, ok := index[key]; ok {
			d := &days[i]
			d.Entries = append(d.Entries, e)
			d.Protein += e.Protein
			d.Carbs += e.Carbs
			d.Fat += e.Fat
			d.Calories = 4*d.Protein + 4*d.Carbs + 9*d.Fat
		} else {
			index[key] = len(days)
			days = append(days, macroDay{
				Date:     e.Date,
				Entries:  []*models.MacroEntry{e},
				Protein:  e.Protein,
				Carbs:    e.Carbs,
				Fat:      e.Fat,
				Calories: 4*e.Protein + 4*e.Carbs + 9*e.Fat,
			})
		}
	}
	return days
}

func (c *MacroController) Index() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	entries, err := Macros.GetAllByUser(userID.(int64))
	if err != nil {
		logs.Error("MacroController.Index: GetAllByUser: %v", err)
	}

	goal, err := MacroGoals.Get(userID.(int64))
	if err != nil {
		logs.Error("MacroController.Index: MacroGoals.Get: %v", err)
	}

	days := groupMacrosByDay(entries)

	distinctFoods, err := Macros.GetDistinctFoods(userID.(int64))
	if err != nil {
		logs.Error("MacroController.Index: GetDistinctFoods: %v", err)
	}
	type foodJSON struct {
		Name          string  `json:"name"`
		ServingWeight float64 `json:"servingWeight"`
		ServingUnit   string  `json:"servingUnit"`
		Protein       float64 `json:"protein"`
		Carbs         float64 `json:"carbs"`
		Fat           float64 `json:"fat"`
	}
	foodList := make([]foodJSON, 0, len(distinctFoods))
	for _, f := range distinctFoods {
		foodList = append(foodList, foodJSON{
			Name:          f.FoodName,
			ServingWeight: f.ServingWeight,
			ServingUnit:   f.ServingUnit,
			Protein:       f.Protein,
			Carbs:         f.Carbs,
			Fat:           f.Fat,
		})
	}
	foodJSON2, _ := json.Marshal(foodList)

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "macros"
	c.Data["Days"] = days
	c.Data["DefaultDate"] = time.Now().Format("2006-01-02")
	c.Data["Goal"] = goal
	c.Data["Summary"] = buildMacroSummary(days, goal)
	c.Data["FoodHistoryJSON"] = template.JS(foodJSON2)
	c.TplName = "macros/index.tpl"
}

func (c *MacroController) SaveGoal() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	protein, _ := strconv.ParseFloat(c.GetString("protein_goal"), 64)
	carbs, _ := strconv.ParseFloat(c.GetString("carbs_goal"), 64)
	fat, _ := strconv.ParseFloat(c.GetString("fat_goal"), 64)

	if _, err := MacroGoals.Upsert(userID.(int64), protein, carbs, fat); err != nil {
		logs.Error("MacroController.SaveGoal: %v", err)
	}

	c.Redirect("/macros", 302)
}

func (c *MacroController) Create() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	dateStr := c.GetString("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.Redirect("/macros", 302)
		return
	}

	foodName := c.GetString("food_name")
	if foodName == "" {
		c.Redirect("/macros", 302)
		return
	}

	servingWeight, _ := strconv.ParseFloat(c.GetString("serving_weight"), 64)
	servingUnit := c.GetString("serving_unit")
	switch servingUnit {
	case "oz", "ml", "fl oz":
	default:
		servingUnit = "g"
	}
	protein, _ := strconv.ParseFloat(c.GetString("protein"), 64)
	carbs, _ := strconv.ParseFloat(c.GetString("carbs"), 64)
	fat, _ := strconv.ParseFloat(c.GetString("fat"), 64)

	if _, err := Macros.Create(userID.(int64), date, foodName, servingWeight, servingUnit, protein, carbs, fat); err != nil {
		logs.Error("MacroController.Create: %v", err)
	}

	c.Redirect("/macros", 302)
}

func (c *MacroController) Update() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Redirect("/macros", 302)
		return
	}

	foodName := c.GetString("food_name")
	if foodName == "" {
		c.Redirect("/macros", 302)
		return
	}

	servingWeight, _ := strconv.ParseFloat(c.GetString("serving_weight"), 64)
	servingUnit := c.GetString("serving_unit")
	switch servingUnit {
	case "oz", "ml", "fl oz":
	default:
		servingUnit = "g"
	}
	protein, _ := strconv.ParseFloat(c.GetString("protein"), 64)
	carbs, _ := strconv.ParseFloat(c.GetString("carbs"), 64)
	fat, _ := strconv.ParseFloat(c.GetString("fat"), 64)

	if _, err := Macros.Update(id, userID.(int64), foodName, servingWeight, servingUnit, protein, carbs, fat); err != nil {
		logs.Error("MacroController.Update: %v", err)
	}

	c.Redirect("/macros", 302)
}

func (c *MacroController) Delete() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.Redirect("/macros", 302)
		return
	}

	if err := Macros.Delete(id, userID.(int64)); err != nil {
		logs.Error("MacroController.Delete: %v", err)
	}

	c.Redirect("/macros", 302)
}
