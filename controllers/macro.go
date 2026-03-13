package controllers

import (
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

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "macros"
	c.Data["Days"] = groupMacrosByDay(entries)
	c.Data["DefaultDate"] = time.Now().Format("2006-01-02")
	c.TplName = "macros/index.tpl"
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
