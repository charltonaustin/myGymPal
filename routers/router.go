package routers

import (
	"myGymPal/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	// Auth
	beego.Router("/register", &controllers.AuthController{}, "get:Register;post:RegisterPost")
	beego.Router("/login", &controllers.AuthController{}, "get:Login;post:LoginPost")

	// App
	beego.Router("/dashboard", &controllers.DashboardController{})
}
