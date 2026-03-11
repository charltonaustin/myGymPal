package controllers_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"myGymPal/models"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

const testUserID = int64(42)
const testProgramID = int64(1)
const testTemplateID = int64(10)
const testSessionID = int64(99)
const testExerciseID = int64(77)
const testSetID = int64(55)

var testPasswordHash string

func init() {
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	if err != nil {
		panic("failed to generate test password hash: " + err.Error())
	}
	testPasswordHash = string(hash)
}

// --- Mock implementations ---

type mockUserRepo struct {
	CreateFn           func(username, password, weightUnit string) (*models.User, error)
	GetByUsernameFn    func(username string) (*models.User, error)
	GetByIDFn          func(id int64) (*models.User, error)
	UpdateWeightUnitFn func(userID int64, unit string) error
	DeleteByUsernameFn func(username string) error
	DeleteByIDFn       func(id int64) error
}

func (m *mockUserRepo) Create(username, password, weightUnit string) (*models.User, error) {
	if m.CreateFn != nil {
		return m.CreateFn(username, password, weightUnit)
	}
	return &models.User{ID: testUserID, Username: username, WeightUnit: weightUnit}, nil
}

func (m *mockUserRepo) GetByUsername(username string) (*models.User, error) {
	if m.GetByUsernameFn != nil {
		return m.GetByUsernameFn(username)
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) GetByID(id int64) (*models.User, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id)
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) UpdateWeightUnit(userID int64, unit string) error {
	if m.UpdateWeightUnitFn != nil {
		return m.UpdateWeightUnitFn(userID, unit)
	}
	return nil
}

func (m *mockUserRepo) DeleteByUsername(username string) error {
	if m.DeleteByUsernameFn != nil {
		return m.DeleteByUsernameFn(username)
	}
	return nil
}

func (m *mockUserRepo) DeleteByID(id int64) error {
	if m.DeleteByIDFn != nil {
		return m.DeleteByIDFn(id)
	}
	return nil
}

type mockProgramRepo struct {
	CreateFn       func(userID int64, name string, startDate time.Time, numPhases, weeksPerPhase, workoutsPerWeek, defaultRepMin, defaultRepMax int) (*models.Program, error)
	GetAllByUserFn func(userID int64) ([]*models.Program, error)
	GetByIDFn      func(id, userID int64) (*models.Program, error)
	DeleteFn       func(id, userID int64) error
}

func (m *mockProgramRepo) Create(userID int64, name string, startDate time.Time, numPhases, weeksPerPhase, workoutsPerWeek, defaultRepMin, defaultRepMax int) (*models.Program, error) {
	if m.CreateFn != nil {
		return m.CreateFn(userID, name, startDate, numPhases, weeksPerPhase, workoutsPerWeek, defaultRepMin, defaultRepMax)
	}
	return &models.Program{ID: testProgramID, UserID: userID, Name: name, NumPhases: numPhases, WeeksPerPhase: weeksPerPhase, WorkoutsPerWeek: workoutsPerWeek}, nil
}

func (m *mockProgramRepo) GetAllByUser(userID int64) ([]*models.Program, error) {
	if m.GetAllByUserFn != nil {
		return m.GetAllByUserFn(userID)
	}
	return nil, nil
}

func (m *mockProgramRepo) GetByID(id, userID int64) (*models.Program, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id, userID)
	}
	return nil, errors.New("not found")
}

func (m *mockProgramRepo) Delete(id, userID int64) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id, userID)
	}
	return nil
}

type mockPhaseRepo struct {
	GetByProgramFn    func(programID int64) ([]*models.Phase, error)
	UpdateRepRangesFn func(programID int64, updates []models.PhaseUpdate) error
}

type mockSessionRepo struct {
	CreateFn         func(programID, userID int64, phaseNumber, weekNumber, workoutNumber int, isDeload bool, date time.Time) (*models.Session, error)
	CountByProgramFn func(programID int64) (int, error)
	GetByIDFn        func(id, userID int64) (*models.Session, error)
	GetByProgramFn   func(programID int64) ([]*models.Session, error)
	DeleteFn         func(id, userID int64) error
}

func (m *mockSessionRepo) Create(programID, userID int64, phaseNumber, weekNumber, workoutNumber int, isDeload bool, date time.Time) (*models.Session, error) {
	if m.CreateFn != nil {
		return m.CreateFn(programID, userID, phaseNumber, weekNumber, workoutNumber, isDeload, date)
	}
	return &models.Session{ID: testSessionID, ProgramID: programID, UserID: userID, PhaseNumber: phaseNumber, WeekNumber: weekNumber, WorkoutNumber: workoutNumber, IsDeload: isDeload}, nil
}

func (m *mockSessionRepo) CountByProgram(programID int64) (int, error) {
	if m.CountByProgramFn != nil {
		return m.CountByProgramFn(programID)
	}
	return 0, nil
}

func (m *mockSessionRepo) GetByID(id, userID int64) (*models.Session, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id, userID)
	}
	return nil, errors.New("not found")
}

func (m *mockSessionRepo) GetByProgram(programID int64) ([]*models.Session, error) {
	if m.GetByProgramFn != nil {
		return m.GetByProgramFn(programID)
	}
	return []*models.Session{}, nil
}

func (m *mockSessionRepo) Delete(id, userID int64) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id, userID)
	}
	return nil
}

type mockTemplateRepo struct {
	CreateFn  func(name, focus string, exercises []models.TemplateExerciseInput) (*models.Template, error)
	UpdateFn  func(id int64, name, focus string, exercises []models.TemplateExerciseInput) (*models.Template, error)
	GetAllFn  func() ([]*models.Template, error)
	GetByIDFn func(id int64) (*models.Template, []*models.TemplateExercise, error)
	DeleteFn  func(id int64) error
}

func (m *mockTemplateRepo) Create(name, focus string, exercises []models.TemplateExerciseInput) (*models.Template, error) {
	if m.CreateFn != nil {
		return m.CreateFn(name, focus, exercises)
	}
	return &models.Template{ID: testTemplateID, Name: name, Focus: focus}, nil
}

func (m *mockTemplateRepo) GetAll() ([]*models.Template, error) {
	if m.GetAllFn != nil {
		return m.GetAllFn()
	}
	return nil, nil
}

func (m *mockTemplateRepo) GetByID(id int64) (*models.Template, []*models.TemplateExercise, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id)
	}
	return nil, nil, errors.New("not found")
}

func (m *mockTemplateRepo) Update(id int64, name, focus string, exercises []models.TemplateExerciseInput) (*models.Template, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(id, name, focus, exercises)
	}
	return &models.Template{ID: id, Name: name, Focus: focus}, nil
}

func (m *mockTemplateRepo) Delete(id int64) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}

func (m *mockPhaseRepo) GetByProgram(programID int64) ([]*models.Phase, error) {
	if m.GetByProgramFn != nil {
		return m.GetByProgramFn(programID)
	}
	return nil, nil
}

func (m *mockPhaseRepo) UpdateRepRanges(programID int64, updates []models.PhaseUpdate) error {
	if m.UpdateRepRangesFn != nil {
		return m.UpdateRepRangesFn(programID, updates)
	}
	return nil
}

type mockSessionExerciseRepo struct {
	CreateFn              func(sessionID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, goalReps int) (*models.SessionExercise, error)
	GetBySessionFn        func(sessionID int64) ([]*models.SessionExerciseView, error)
	GetByIDFn             func(exerciseID int64) (*models.SessionExercise, error)
	LogSetFn              func(exerciseID int64, setNumber int, actualWeight float64, weightUnit string, actualReps int) (*models.SessionSet, error)
	CountSetsByExerciseFn func(exerciseID int64) (int, error)
	DeleteSetFn           func(setID int64) error
}

func (m *mockSessionExerciseRepo) Create(sessionID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, goalReps int) (*models.SessionExercise, error) {
	if m.CreateFn != nil {
		return m.CreateFn(sessionID, name, isBodyweight, goalWeight, weightUnit, goalReps)
	}
	return &models.SessionExercise{ID: testExerciseID, SessionID: sessionID, Name: name, IsBodyweight: isBodyweight, GoalWeight: goalWeight, WeightUnit: weightUnit, GoalReps: goalReps}, nil
}

func (m *mockSessionExerciseRepo) GetBySession(sessionID int64) ([]*models.SessionExerciseView, error) {
	if m.GetBySessionFn != nil {
		return m.GetBySessionFn(sessionID)
	}
	return []*models.SessionExerciseView{}, nil
}

func (m *mockSessionExerciseRepo) GetByID(exerciseID int64) (*models.SessionExercise, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(exerciseID)
	}
	return nil, errors.New("not found")
}

func (m *mockSessionExerciseRepo) LogSet(exerciseID int64, setNumber int, actualWeight float64, weightUnit string, actualReps int) (*models.SessionSet, error) {
	if m.LogSetFn != nil {
		return m.LogSetFn(exerciseID, setNumber, actualWeight, weightUnit, actualReps)
	}
	return &models.SessionSet{ID: testSetID, SessionExerciseID: exerciseID, SetNumber: setNumber, ActualWeight: actualWeight, WeightUnit: weightUnit, ActualReps: actualReps}, nil
}

func (m *mockSessionExerciseRepo) CountSetsByExercise(exerciseID int64) (int, error) {
	if m.CountSetsByExerciseFn != nil {
		return m.CountSetsByExerciseFn(exerciseID)
	}
	return 0, nil
}

func (m *mockSessionExerciseRepo) DeleteSet(setID int64) error {
	if m.DeleteSetFn != nil {
		return m.DeleteSetFn(setID)
	}
	return nil
}

type mockExerciseRepo struct {
	CreateFn           func(userID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string) (*models.Exercise, error)
	GetAllByUserFn     func(userID int64) ([]*models.Exercise, error)
	GetByIDFn          func(id, userID int64) (*models.Exercise, error)
	GetByNameFn        func(userID int64, name string) (*models.Exercise, error)
	UpdateFn           func(id, userID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string) (*models.Exercise, error)
	UpdateGoalWeightFn func(id int64, goalWeight float64) error
	DeleteFn           func(id, userID int64) error
}

func (m *mockExerciseRepo) Create(userID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string) (*models.Exercise, error) {
	if m.CreateFn != nil {
		return m.CreateFn(userID, name, isBodyweight, goalWeight, weightUnit)
	}
	return &models.Exercise{ID: 1, UserID: userID, Name: name, IsBodyweight: isBodyweight, GoalWeight: goalWeight, WeightUnit: weightUnit}, nil
}

func (m *mockExerciseRepo) GetAllByUser(userID int64) ([]*models.Exercise, error) {
	if m.GetAllByUserFn != nil {
		return m.GetAllByUserFn(userID)
	}
	return []*models.Exercise{}, nil
}

func (m *mockExerciseRepo) GetByID(id, userID int64) (*models.Exercise, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id, userID)
	}
	return nil, errors.New("not found")
}

func (m *mockExerciseRepo) GetByName(userID int64, name string) (*models.Exercise, error) {
	if m.GetByNameFn != nil {
		return m.GetByNameFn(userID, name)
	}
	return nil, errors.New("not found")
}

func (m *mockExerciseRepo) Update(id, userID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string) (*models.Exercise, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(id, userID, name, isBodyweight, goalWeight, weightUnit)
	}
	return &models.Exercise{ID: id, UserID: userID, Name: name, IsBodyweight: isBodyweight, GoalWeight: goalWeight, WeightUnit: weightUnit}, nil
}

func (m *mockExerciseRepo) UpdateGoalWeight(id int64, goalWeight float64) error {
	if m.UpdateGoalWeightFn != nil {
		return m.UpdateGoalWeightFn(id, goalWeight)
	}
	return nil
}

func (m *mockExerciseRepo) Delete(id, userID int64) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id, userID)
	}
	return nil
}

// --- Global mock instances ---

var (
	mockUsers            = &mockUserRepo{}
	mockPrograms         = &mockProgramRepo{}
	mockPhases           = &mockPhaseRepo{}
	mockSessions         = &mockSessionRepo{}
	mockTemplates        = &mockTemplateRepo{}
	mockSessionExercises = &mockSessionExerciseRepo{}
	mockExercises        = &mockExerciseRepo{}
)

func resetMocks() {
	*mockUsers = mockUserRepo{}
	*mockPrograms = mockProgramRepo{}
	*mockPhases = mockPhaseRepo{}
	*mockSessions = mockSessionRepo{}
	*mockTemplates = mockTemplateRepo{}
	*mockSessionExercises = mockSessionExerciseRepo{}
	*mockExercises = mockExerciseRepo{}
	lastProgramCreate = struct {
		name            string
		numPhases       int
		weeksPerPhase   int
		workoutsPerWeek int
	}{}
	lastPhaseUpdates = nil
	lastTemplateCreate = struct {
		name         string
		focus        string
		numExercises int
	}{}
	lastSessionCreate = struct {
		phaseNumber   int
		weekNumber    int
		workoutNumber int
		isDeload      bool
	}{}
	lastLogSet = struct {
		exerciseID   int64
		setNumber    int
		actualWeight float64
		weightUnit   string
		actualReps   int
	}{}
	sessionExerciseCreateNames = nil
}

// setGetByUsernameReturnsUser makes mockUsers return a test user for any username lookup.
func setGetByUsernameReturnsUser(weightUnit string) {
	mockUsers.GetByUsernameFn = func(username string) (*models.User, error) {
		return newTestUser(username, weightUnit), nil
	}
}

// setCreateFnError makes mockUsers.Create return an error.
func setCreateFnError(err error) {
	mockUsers.CreateFn = func(username, password, weightUnit string) (*models.User, error) {
		return nil, err
	}
}

// --- Program mock helpers ---

// setProgramsGetAllEmpty makes GetAllByUser return an empty slice.
func setProgramsGetAllEmpty() {
	mockPrograms.GetAllByUserFn = func(userID int64) ([]*models.Program, error) {
		return []*models.Program{}, nil
	}
}

// setProgramsGetAllWithOne makes GetAllByUser return a single program.
func setProgramsGetAllWithOne(id int64, name string, numPhases, weeksPerPhase int) {
	mockPrograms.GetAllByUserFn = func(userID int64) ([]*models.Program, error) {
		return []*models.Program{
			{ID: id, UserID: userID, Name: name, NumPhases: numPhases, WeeksPerPhase: weeksPerPhase},
		}, nil
	}
}

// setProgramGetByIDError makes GetByID return an error.
func setProgramGetByIDError(err error) {
	mockPrograms.GetByIDFn = func(id, userID int64) (*models.Program, error) {
		return nil, err
	}
}

// setProgramGetByID makes GetByID return a program with the given name and phase count.
func setProgramGetByID(name string, numPhases int) {
	mockPrograms.GetByIDFn = func(id, userID int64) (*models.Program, error) {
		return &models.Program{ID: id, UserID: userID, Name: name, NumPhases: numPhases}, nil
	}
}

// setPhasesGetByProgram makes GetByProgram return n phases with default rep ranges (10–12).
func setPhasesGetByProgram(count int) {
	mockPhases.GetByProgramFn = func(programID int64) ([]*models.Phase, error) {
		phases := make([]*models.Phase, count)
		for i := range phases {
			phases[i] = &models.Phase{
				ID:          int64(i + 1),
				ProgramID:   programID,
				PhaseNumber: i + 1,
				RepMin:      10,
				RepMax:      12,
			}
		}
		return phases, nil
	}
}

// lastProgramCreate holds args captured by captureProgramCreate.
var lastProgramCreate struct {
	name            string
	numPhases       int
	weeksPerPhase   int
	workoutsPerWeek int
}

// captureProgramCreate makes CreateFn capture the call args and return a valid program.
func captureProgramCreate() {
	mockPrograms.CreateFn = func(userID int64, name string, startDate time.Time, numPhases, weeksPerPhase, workoutsPerWeek, defaultRepMin, defaultRepMax int) (*models.Program, error) {
		lastProgramCreate.name = name
		lastProgramCreate.numPhases = numPhases
		lastProgramCreate.weeksPerPhase = weeksPerPhase
		lastProgramCreate.workoutsPerWeek = workoutsPerWeek
		return &models.Program{ID: testProgramID, UserID: userID, Name: name, NumPhases: numPhases, WeeksPerPhase: weeksPerPhase, WorkoutsPerWeek: workoutsPerWeek}, nil
	}
}

// lastPhaseUpdates holds args captured by capturePhasesUpdateRepRanges.
var lastPhaseUpdates []models.PhaseUpdate

// capturePhasesUpdateRepRanges makes UpdateRepRangesFn capture the call args.
func capturePhasesUpdateRepRanges() {
	mockPhases.UpdateRepRangesFn = func(programID int64, updates []models.PhaseUpdate) error {
		lastPhaseUpdates = updates
		return nil
	}
}

// getLastPhaseUpdate returns the RepMin and RepMax for the i-th captured update.
func getLastPhaseUpdate(i int) (repMin, repMax int) {
	return lastPhaseUpdates[i].RepMin, lastPhaseUpdates[i].RepMax
}

// setPhasesUpdateRepRangesError makes UpdateRepRangesFn return an error.
func setPhasesUpdateRepRangesError(err error) {
	mockPhases.UpdateRepRangesFn = func(programID int64, updates []models.PhaseUpdate) error {
		return err
	}
}

// --- Session mock helpers ---

// lastSessionCreate holds args captured by captureSessionCreate.
var lastSessionCreate struct {
	phaseNumber   int
	weekNumber    int
	workoutNumber int
	isDeload      bool
}

// captureSessionCreate makes CreateFn capture the call args and return a valid session.
func captureSessionCreate() {
	mockSessions.CreateFn = func(programID, userID int64, phaseNumber, weekNumber, workoutNumber int, isDeload bool, date time.Time) (*models.Session, error) {
		lastSessionCreate.phaseNumber = phaseNumber
		lastSessionCreate.weekNumber = weekNumber
		lastSessionCreate.workoutNumber = workoutNumber
		lastSessionCreate.isDeload = isDeload
		return &models.Session{ID: testSessionID, ProgramID: programID, UserID: userID, PhaseNumber: phaseNumber, WeekNumber: weekNumber, WorkoutNumber: workoutNumber, IsDeload: isDeload}, nil
	}
}

// setSessionsGetByProgram makes GetByProgram return n sessions.
func setSessionsGetByProgram(count int) {
	mockSessions.GetByProgramFn = func(programID int64) ([]*models.Session, error) {
		sessions := make([]*models.Session, count)
		for i := range sessions {
			sessions[i] = &models.Session{
				ID:            int64(i + 1),
				ProgramID:     programID,
				PhaseNumber:   1,
				WeekNumber:    i + 1,
				WorkoutNumber: i + 1,
			}
		}
		return sessions, nil
	}
}

// setSessionCountByProgram makes CountByProgram return a fixed count.
func setSessionCountByProgram(count int) {
	mockSessions.CountByProgramFn = func(programID int64) (int, error) {
		return count, nil
	}
}

// setSessionGetByIDError makes GetByID return an error.
func setSessionGetByIDError(err error) {
	mockSessions.GetByIDFn = func(id, userID int64) (*models.Session, error) {
		return nil, err
	}
}

// setSessionGetByID makes GetByID return a session with the given values.
func setSessionGetByID(phaseNumber, weekNumber, workoutNumber int, isDeload bool) {
	mockSessions.GetByIDFn = func(id, userID int64) (*models.Session, error) {
		return &models.Session{
			ID:            id,
			UserID:        userID,
			ProgramID:     testProgramID,
			PhaseNumber:   phaseNumber,
			WeekNumber:    weekNumber,
			WorkoutNumber: workoutNumber,
			IsDeload:      isDeload,
		}, nil
	}
}

// setProgramGetByIDWithDates makes GetByID return a program with StartDate, WeeksPerPhase,
// and WorkoutsPerWeek (defaulting to 4) set.
func setProgramGetByIDWithDates(name string, numPhases, weeksPerPhase int, startDate time.Time) {
	mockPrograms.GetByIDFn = func(id, userID int64) (*models.Program, error) {
		return &models.Program{
			ID:              id,
			UserID:          userID,
			Name:            name,
			NumPhases:       numPhases,
			WeeksPerPhase:   weeksPerPhase,
			WorkoutsPerWeek: 4,
			StartDate:       startDate,
		}, nil
	}
}

// --- Template mock helpers ---

// setTemplatesGetAllEmpty makes GetAll return an empty slice.
func setTemplatesGetAllEmpty() {
	mockTemplates.GetAllFn = func() ([]*models.Template, error) {
		return []*models.Template{}, nil
	}
}

// setTemplatesGetAllWithOne makes GetAll return a single template.
func setTemplatesGetAllWithOne(id int64, name, focus string) {
	mockTemplates.GetAllFn = func() ([]*models.Template, error) {
		return []*models.Template{
			{ID: id, Name: name, Focus: focus},
		}, nil
	}
}

// setTemplateGetByIDError makes GetByID return an error.
func setTemplateGetByIDError(err error) {
	mockTemplates.GetByIDFn = func(id int64) (*models.Template, []*models.TemplateExercise, error) {
		return nil, nil, err
	}
}

// setTemplateGetByID makes GetByID return a template with n exercises.
func setTemplateGetByID(id int64, name, focus string, numExercises int) {
	mockTemplates.GetByIDFn = func(_ int64) (*models.Template, []*models.TemplateExercise, error) {
		exercises := make([]*models.TemplateExercise, numExercises)
		for i := range exercises {
			exercises[i] = &models.TemplateExercise{
				ID:         int64(i + 1),
				TemplateID: id,
				Name:       "Exercise " + fmt.Sprintf("%d", i+1),
				SortOrder:  i,
			}
		}
		return &models.Template{ID: id, Name: name, Focus: focus}, exercises, nil
	}
}

// lastTemplateCreate holds args captured by captureTemplateCreate.
var lastTemplateCreate struct {
	name         string
	focus        string
	numExercises int
}

// captureTemplateCreate makes CreateFn capture the call args and return a valid template.
func captureTemplateCreate() {
	mockTemplates.CreateFn = func(name, focus string, exercises []models.TemplateExerciseInput) (*models.Template, error) {
		lastTemplateCreate.name = name
		lastTemplateCreate.focus = focus
		lastTemplateCreate.numExercises = len(exercises)
		return &models.Template{ID: testTemplateID, Name: name, Focus: focus}, nil
	}
}

// setTemplateCreateError makes CreateFn return an error.
func setTemplateCreateError(err error) {
	mockTemplates.CreateFn = func(name, focus string, exercises []models.TemplateExerciseInput) (*models.Template, error) {
		return nil, err
	}
}

// --- SessionExercise mock helpers ---

// lastLogSet holds args captured by captureLogSet.
var lastLogSet struct {
	exerciseID   int64
	setNumber    int
	actualWeight float64
	weightUnit   string
	actualReps   int
}

// sessionExerciseCreateNames captures exercise names from captureSessionExerciseCreates.
var sessionExerciseCreateNames []string

// captureSessionExerciseCreates records the name of each exercise created.
func captureSessionExerciseCreates() {
	sessionExerciseCreateNames = nil
	mockSessionExercises.CreateFn = func(sessionID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, goalReps int) (*models.SessionExercise, error) {
		sessionExerciseCreateNames = append(sessionExerciseCreateNames, name)
		return &models.SessionExercise{ID: testExerciseID, SessionID: sessionID, Name: name, IsBodyweight: isBodyweight, GoalWeight: goalWeight, WeightUnit: weightUnit, GoalReps: goalReps}, nil
	}
}

// setSessionExerciseGetBySessionWithOne makes GetBySession return a single exercise with no sets.
func setSessionExerciseGetBySessionWithOne(name, weightUnit string) {
	mockSessionExercises.GetBySessionFn = func(sessionID int64) ([]*models.SessionExerciseView, error) {
		ex := &models.SessionExercise{ID: testExerciseID, SessionID: sessionID, Name: name, WeightUnit: weightUnit}
		return []*models.SessionExerciseView{{Exercise: ex, Sets: nil}}, nil
	}
}

// setSessionExerciseGetByID makes GetByID return an exercise belonging to the given session.
func setSessionExerciseGetByID(sessionID int64) {
	mockSessionExercises.GetByIDFn = func(exerciseID int64) (*models.SessionExercise, error) {
		return &models.SessionExercise{ID: exerciseID, SessionID: sessionID, Name: "Bench Press", WeightUnit: "lb"}, nil
	}
}

// setSessionExerciseGetByIDError makes GetByID return an error.
func setSessionExerciseGetByIDError(err error) {
	mockSessionExercises.GetByIDFn = func(exerciseID int64) (*models.SessionExercise, error) {
		return nil, err
	}
}

// captureLogSet makes LogSetFn capture the call args and return a valid set.
func captureLogSet() {
	mockSessionExercises.LogSetFn = func(exerciseID int64, setNumber int, actualWeight float64, weightUnit string, actualReps int) (*models.SessionSet, error) {
		lastLogSet.exerciseID = exerciseID
		lastLogSet.setNumber = setNumber
		lastLogSet.actualWeight = actualWeight
		lastLogSet.weightUnit = weightUnit
		lastLogSet.actualReps = actualReps
		return &models.SessionSet{ID: testSetID, SessionExerciseID: exerciseID, SetNumber: setNumber, ActualWeight: actualWeight, WeightUnit: weightUnit, ActualReps: actualReps}, nil
	}
}

// setLogSetCountByExercise makes CountSetsByExercise return a fixed count.
func setLogSetCountByExercise(count int) {
	mockSessionExercises.CountSetsByExerciseFn = func(exerciseID int64) (int, error) {
		return count, nil
	}
}

// --- Test helpers ---

func newTestUser(username, weightUnit string) *models.User {
	return &models.User{
		ID:           testUserID,
		Username:     username,
		PasswordHash: testPasswordHash,
		WeightUnit:   weightUnit,
	}
}

// loginAs configures the mock user repo so the given user can log in,
// performs POST /login, and returns the session cookies.
func loginAs(t *testing.T, username, weightUnit string) []*http.Cookie {
	t.Helper()
	user := newTestUser(username, weightUnit)
	mockUsers.GetByUsernameFn = func(u string) (*models.User, error) {
		if u == username {
			return user, nil
		}
		return nil, errors.New("not found")
	}
	mockUsers.GetByIDFn = func(id int64) (*models.User, error) {
		if id == testUserID {
			return user, nil
		}
		return nil, errors.New("not found")
	}
	w := postForm("/login", url.Values{
		"username": {username},
		"password": {"password123"},
	}, nil)
	require.Equal(t, http.StatusFound, w.Code)
	return w.Result().Cookies()
}

func postForm(path string, data url.Values, cookies []*http.Cookie) *httptest.ResponseRecorder {
	r, _ := http.NewRequest("POST", path, strings.NewReader(data.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, c := range cookies {
		r.AddCookie(c)
	}
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	return w
}

func getPath(path string, cookies []*http.Cookie) *httptest.ResponseRecorder {
	r, _ := http.NewRequest("GET", path, nil)
	for _, c := range cookies {
		r.AddCookie(c)
	}
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	return w
}
