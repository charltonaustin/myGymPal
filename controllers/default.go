package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	if c.GetSession("user_id") != nil {
		c.Redirect("/dashboard", 302)
		return
	}
	c.Data["ActivePage"] = "home"
	c.TplName = "index.tpl"
}
