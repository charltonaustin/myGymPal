package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["json"] = map[string]string{"status": "ok", "app": "myGymPal"}
	c.ServeJSON()
}
