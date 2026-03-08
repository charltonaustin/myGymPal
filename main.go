package main

import (
	"log"

	_ "myGymPal/routers"
	"myGymPal/models"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	if err := models.Init(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	beego.Run()
}
