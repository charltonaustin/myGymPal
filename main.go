package main

import (
	"log"

	"myGymPal/models"
	"myGymPal/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	if err := models.Init(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	routers.Register()
	beego.Run()
}
