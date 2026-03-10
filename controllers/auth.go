package controllers

import (
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

type AuthController struct {
	beego.Controller
}

func (c *AuthController) Register() {
	c.TplName = "auth/register.tpl"
}

func (c *AuthController) Login() {
	c.TplName = "auth/login.tpl"
}

func (c *AuthController) Logout() {
	if err := c.DestroySession(); err != nil {
		c.Redirect("/error", 302)
		return
	}
	c.Redirect("/login", 302)
}

func (c *AuthController) LoginPost() {
	username := strings.TrimSpace(c.GetString("username"))
	password := c.GetString("password")

	user, err := Users.GetByUsername(username)
	if err != nil || !user.CheckPassword(password) {
		c.Data["Error"] = "Invalid username or password."
		c.Data["Username"] = username
		c.TplName = "auth/login.tpl"
		return
	}

	err = c.SetSession("user_id", user.ID)
	if err != nil {
		c.Redirect("/error", 302)
		return
	}
	err = c.SetSession("username", user.Username)
	if err != nil {
		c.Redirect("/error", 302)
		return
	}
	c.Redirect("/dashboard", 302)
}

func (c *AuthController) RegisterPost() {
	username := strings.TrimSpace(c.GetString("username"))
	password := c.GetString("password")
	confirm := c.GetString("confirm_password")
	weightUnit := "lb"

	if username == "" || password == "" {
		c.Data["Error"] = "Username and password are required."
		c.Data["Username"] = username
		c.TplName = "auth/register.tpl"
		return
	}

	if len(password) < 8 {
		c.Data["Error"] = "Password must be at least 8 characters."
		c.Data["Username"] = username
		c.TplName = "auth/register.tpl"
		return
	}

	if password != confirm {
		c.Data["Error"] = "Passwords do not match."
		c.Data["Username"] = username
		c.TplName = "auth/register.tpl"
		return
	}

	if _, err := Users.Create(username, password, weightUnit); err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			c.Data["Error"] = "That username is already taken."
		} else {
			c.Data["Error"] = "Something went wrong. Please try again."
		}
		c.Data["Username"] = username
		c.TplName = "auth/register.tpl"
		return
	}

	c.Redirect("/login", 302)
}
