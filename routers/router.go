package routers

import (
	"myGymPal/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func Register() {
	beego.Router("/", &controllers.MainController{})

	// Auth
	beego.Router("/register", &controllers.AuthController{}, "get:Register;post:RegisterPost")
	beego.Router("/login", &controllers.AuthController{}, "get:Login;post:LoginPost")
	beego.Router("/logout", &controllers.AuthController{}, "get:Logout")

	// App
	beego.Router("/dashboard", &controllers.DashboardController{})
}
