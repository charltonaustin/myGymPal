package controllers

import (
	"os"

	beego "github.com/beego/beego/v2/server/web"
)

type PWAController struct {
	beego.Controller
}

func (c *PWAController) ServiceWorker() {
	content, err := os.ReadFile("static/sw.js")
	if err != nil {
		c.Ctx.Output.SetStatus(404)
		return
	}
	c.Ctx.Output.Header("Content-Type", "application/javascript; charset=utf-8")
	c.Ctx.Output.Header("Cache-Control", "no-cache")
	c.Ctx.Output.Header("Service-Worker-Allowed", "/")
	c.Ctx.Output.Body(content)
}

func (c *PWAController) Manifest() {
	content, err := os.ReadFile("static/manifest.json")
	if err != nil {
		c.Ctx.Output.SetStatus(404)
		return
	}
	c.Ctx.Output.Header("Content-Type", "application/manifest+json")
	c.Ctx.Output.Body(content)
}

func (c *PWAController) Offline() {
	c.TplName = "offline.tpl"
}
