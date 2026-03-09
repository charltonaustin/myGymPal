package controllers

import beego "github.com/beego/beego/v2/server/web"

type DashboardController struct {
	beego.Controller
}

func (c *DashboardController) Get() {
	username := c.GetSession("user_id")
	if username == nil {
		c.Redirect("/login", 302)
		return
	}
	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "dashboard"
	c.Data["Username"] = c.GetSession("username")
	c.TplName = "dashboard.tpl"
}
