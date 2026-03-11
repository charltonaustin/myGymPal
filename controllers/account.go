package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type AccountController struct {
	beego.Controller
}

func (c *AccountController) Settings() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	user, err := Users.GetByID(userID.(int64))
	if err != nil {
		c.Redirect("/error", 302)
		return
	}

	flash := beego.ReadFromRequest(&c.Controller)
	if saved, ok := flash.Data["success"]; ok {
		c.Data["Success"] = saved
	}

	c.Data["LoggedIn"] = true
	c.Data["ActivePage"] = "settings"
	c.Data["WeightUnit"] = user.WeightUnit
	c.TplName = "settings.tpl"
}

func (c *AccountController) DeleteAccount() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	if err := Users.DeleteByID(userID.(int64)); err != nil {
		c.Data["LoggedIn"] = true
		c.Data["ActivePage"] = "settings"
		c.Data["Error"] = "Failed to delete account. Please try again."
		user, _ := Users.GetByID(userID.(int64))
		if user != nil {
			c.Data["WeightUnit"] = user.WeightUnit
		}
		c.TplName = "settings.tpl"
		return
	}

	c.DestroySession()
	c.Redirect("/login", 302)
}

func (c *AccountController) SettingsPost() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	unit := c.GetString("weight_unit")
	if err := Users.UpdateWeightUnit(userID.(int64), unit); err != nil {
		c.Data["LoggedIn"] = true
		c.Data["ActivePage"] = "settings"
		c.Data["Error"] = "Invalid weight unit."
		c.Data["WeightUnit"] = unit
		c.TplName = "settings.tpl"
		return
	}

	flash := beego.NewFlash()
	flash.Success("Settings saved.")
	flash.Store(&c.Controller)
	c.Redirect("/settings", 302)
}
