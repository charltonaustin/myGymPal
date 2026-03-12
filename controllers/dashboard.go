package controllers

import beego "github.com/beego/beego/v2/server/web"

type DashboardController struct {
	beego.Controller
}

func (c *DashboardController) Get() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}
	recent, _ := Sessions.GetRecentByUser(userID.(int64), 10)
	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "dashboard"
	c.Data["Username"] = c.GetSession("username")
	c.Data["RecentSessions"] = recent
	c.TplName = "dashboard.tpl"
}
