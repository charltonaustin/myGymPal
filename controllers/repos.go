package controllers

import "myGymPal/models"

// Repository dependencies injected by routers.Register before the server starts.
// Controllers reference these package-level vars; Beego instantiates controllers
// via reflect.New so field-level injection is not available.
var (
	Users            models.UserRepository
	Programs         models.ProgramRepository
	Phases           models.PhaseRepository
	Templates        models.TemplateRepository
	Sessions         models.SessionRepository
	SessionExercises models.SessionExerciseRepository
	Exercises        models.ExerciseRepository
)
