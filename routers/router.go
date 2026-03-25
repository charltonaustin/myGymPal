package routers

import (
	"fmt"
	"myGymPal/controllers"
	"myGymPal/models"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.AddFuncMap("fmtDuration", func(secs int) string {
		h := secs / 3600
		m := (secs % 3600) / 60
		s := secs % 60
		if h > 0 {
			return fmt.Sprintf("%d:%02d:%02d", h, m, s)
		}
		return fmt.Sprintf("%d:%02d", m, s)
	})
	beego.AddFuncMap("restMinutes", func(secs int) int { return secs / 60 })
	beego.AddFuncMap("restSecs", func(secs int) int { return secs % 60 })
}

func Register() {
	controllers.Users = models.NewUserRepository()
	controllers.Programs = models.NewProgramRepository()
	controllers.Phases = models.NewPhaseRepository()
	controllers.Templates = models.NewTemplateRepository()
	controllers.Sessions = models.NewSessionRepository()
	controllers.SessionExercises = models.NewSessionExerciseRepository()
	controllers.Exercises = models.NewExerciseRepository()
	controllers.BodyWeights = models.NewBodyWeightRepository()
	controllers.Macros = models.NewMacroRepository()
	controllers.MacroGoals = models.NewMacroGoalRepository()

	// PWA assets (must be served from root path for correct service worker scope)
	beego.Router("/sw.js", &controllers.PWAController{}, "get:ServiceWorker")
	beego.Router("/manifest.json", &controllers.PWAController{}, "get:Manifest")
	beego.Router("/offline", &controllers.PWAController{}, "get:Offline")

	beego.Router("/", &controllers.MainController{})
	beego.Router("/example", &controllers.MainController{}, "get:Example")

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
	beego.Router("/sessions/:id/exercises/reorder", &controllers.SessionController{}, "post:ReorderExercises")
	beego.Router("/sessions/:id/exercises/:eid/delete", &controllers.SessionController{}, "post:DeleteExercise")
	beego.Router("/sessions/:id/exercises/:eid/sets", &controllers.SessionController{}, "post:LogSet")
	beego.Router("/sessions/:id/exercises/:eid/sets/:sid/delete", &controllers.SessionController{}, "post:DeleteSet")
	beego.Router("/sessions/:id/cardio", &controllers.SessionController{}, "post:AddCardioActivity")
	beego.Router("/sessions/:id/exercises/:eid/cardio", &controllers.SessionController{}, "post:LogCardio")
	beego.Router("/sessions/:id/exercises/:eid/cardio/:lid/delete", &controllers.SessionController{}, "post:DeleteCardioLog")

	// Macros
	beego.Router("/macros", &controllers.MacroController{}, "get:Index;post:Create")
	beego.Router("/macros/goals", &controllers.MacroController{}, "post:SaveGoal")
	beego.Router("/macros/:id", &controllers.MacroController{}, "post:Update")
	beego.Router("/macros/:id/delete", &controllers.MacroController{}, "post:Delete")

	// Weight
	beego.Router("/weight", &controllers.WeightController{}, "get:Index;post:Create")
	beego.Router("/weight/:id", &controllers.WeightController{}, "post:Update")
	beego.Router("/weight/:id/delete", &controllers.WeightController{}, "post:Delete")

	// Exercises
	beego.Router("/exercises", &controllers.ExerciseController{}, "get:Index")
	beego.Router("/exercises/new", &controllers.ExerciseController{}, "get:New;post:Create")
	beego.Router("/exercises/goal-weight", &controllers.ExerciseController{}, "post:UpdateGoalWeightJSON")
	beego.Router("/exercises/goal-reps", &controllers.ExerciseController{}, "post:UpdateGoalRepsJSON")
	beego.Router("/exercises/goal-seconds", &controllers.ExerciseController{}, "post:UpdateGoalSecondsJSON")
	beego.Router("/exercises/:id/edit", &controllers.ExerciseController{}, "get:Edit;post:Update")
	beego.Router("/exercises/:id/delete", &controllers.ExerciseController{}, "post:Delete")

	// Templates
	beego.Router("/templates", &controllers.TemplateController{}, "get:Index")
	beego.Router("/templates/new", &controllers.TemplateController{}, "get:New;post:Create")
	beego.Router("/templates/:id", &controllers.TemplateController{}, "get:Show;post:Update")
	beego.Router("/templates/:id/edit", &controllers.TemplateController{}, "get:Edit")
	beego.Router("/templates/:id/delete", &controllers.TemplateController{}, "post:Delete")
}
