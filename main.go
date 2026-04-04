package main

import (
	"log"

	"myGymPal/models"
	"myGymPal/routers"

	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/beego/beego/v2/server/web/session/postgres"
)

func main() {
	if err := models.Init(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	beego.BConfig.WebConfig.Session.SessionProviderConfig = models.SessionProviderDSN()

	beego.SetStaticPath("/static", "static")
	routers.Register()
	beego.Run()
}
