package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type DashboardController struct {
	beego.Controller
}

func (c *DashboardController) Get() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	recent, err := Sessions.GetRecentByUser(userID.(int64), 10)
	if err != nil {
		logs.Error("DashboardController.Get: GetRecentByUser: %v", err)
	}

	weightEntries, err := BodyWeights.GetAllByUser(userID.(int64))
	if err != nil {
		logs.Error("DashboardController.Get: GetAllByUser (weight): %v", err)
	}

	macroEntries, err := Macros.GetAllByUser(userID.(int64))
	if err != nil {
		logs.Error("DashboardController.Get: GetAllByUser (macros): %v", err)
	}
	macroDays := groupMacrosByDay(macroEntries)
	macroGoal, err := MacroGoals.Get(userID.(int64))
	if err != nil {
		logs.Error("DashboardController.Get: MacroGoals.Get: %v", err)
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "dashboard"
	c.Data["Username"] = c.GetSession("username")
	c.Data["RecentSessions"] = recent
	c.Data["WeightAvg"] = computeWeightAverage(weightEntries)
	c.Data["MacroSummary"] = buildMacroSummary(macroDays, macroGoal)
	c.TplName = "dashboard.tpl"
}
