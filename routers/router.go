package routers

import (
	"myGymPal/controllers"
	"myGymPal/models"

	beego "github.com/beego/beego/v2/server/web"
)

func Register() {
	controllers.Users = models.NewUserRepository()
	controllers.Programs = models.NewProgramRepository()
	controllers.Phases = models.NewPhaseRepository()
	controllers.Templates = models.NewTemplateRepository()
	controllers.Sessions = models.NewSessionRepository()
	controllers.SessionExercises = models.NewSessionExerciseRepository()

	// PWA assets (must be served from root path for correct service worker scope)
	beego.Router("/sw.js", &controllers.PWAController{}, "get:ServiceWorker")
	beego.Router("/manifest.json", &controllers.PWAController{}, "get:Manifest")
	beego.Router("/offline", &controllers.PWAController{}, "get:Offline")

	beego.Router("/", &controllers.MainController{})

	// Auth
	beego.Router("/register", &controllers.AuthController{}, "get:Register;post:RegisterPost")
	beego.Router("/login", &controllers.AuthController{}, "get:Login;post:LoginPost")
	beego.Router("/logout", &controllers.AuthController{}, "get:Logout")

	// App
	beego.Router("/dashboard", &controllers.DashboardController{})
	beego.Router("/settings", &controllers.AccountController{}, "get:Settings;post:SettingsPost")
	beego.Router("/account/delete", &controllers.AccountController{}, "post:DeleteAccount")
	beego.Router("/error", &controllers.ErrorController{})

	// Programs
	beego.Router("/programs", &controllers.ProgramController{}, "get:Index;post:Create")
	beego.Router("/programs/new", &controllers.ProgramController{}, "get:New")
	beego.Router("/programs/:id", &controllers.ProgramController{}, "get:Show;post:UpdatePhases")
	beego.Router("/programs/:id/delete", &controllers.ProgramController{}, "post:Delete")

	// Sessions
	beego.Router("/programs/:id/sessions/new", &controllers.SessionController{}, "get:New")
	beego.Router("/programs/:id/sessions", &controllers.SessionController{}, "post:Create")
	beego.Router("/sessions/:id", &controllers.SessionController{}, "get:Show")
	beego.Router("/sessions/:id/delete", &controllers.SessionController{}, "post:Delete")
	beego.Router("/sessions/:id/exercises", &controllers.SessionController{}, "post:AddExercise")
	beego.Router("/sessions/:id/exercises/:eid/sets", &controllers.SessionController{}, "post:LogSet")

	// Templates
	beego.Router("/templates", &controllers.TemplateController{}, "get:Index")
	beego.Router("/templates/new", &controllers.TemplateController{}, "get:New;post:Create")
	beego.Router("/templates/:id", &controllers.TemplateController{}, "get:Show;post:Update")
	beego.Router("/templates/:id/edit", &controllers.TemplateController{}, "get:Edit")
	beego.Router("/templates/:id/delete", &controllers.TemplateController{}, "post:Delete")
}
