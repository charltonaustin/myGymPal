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
	beego.Router("/settings", &controllers.AccountController{}, "get:Settings;post:SettingsPost")
	beego.Router("/error", &controllers.ErrorController{})

	// Programs
	beego.Router("/programs", &controllers.ProgramController{}, "get:Index;post:Create")
	beego.Router("/programs/new", &controllers.ProgramController{}, "get:New")
}
