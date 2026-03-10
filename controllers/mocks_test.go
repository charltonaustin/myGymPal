package controllers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"myGymPal/models"
)

const testUserID = int64(42)
const testProgramID = int64(1)

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

type mockProgramRepo struct {
	CreateFn       func(userID int64, name string, startDate time.Time, numPhases, weeksPerPhase, defaultRepMin, defaultRepMax int) (*models.Program, error)
	GetAllByUserFn func(userID int64) ([]*models.Program, error)
	GetByIDFn      func(id, userID int64) (*models.Program, error)
	DeleteFn       func(id, userID int64) error
}

func (m *mockProgramRepo) Create(userID int64, name string, startDate time.Time, numPhases, weeksPerPhase, defaultRepMin, defaultRepMax int) (*models.Program, error) {
	if m.CreateFn != nil {
		return m.CreateFn(userID, name, startDate, numPhases, weeksPerPhase, defaultRepMin, defaultRepMax)
	}
	return &models.Program{ID: testProgramID, UserID: userID, Name: name, NumPhases: numPhases, WeeksPerPhase: weeksPerPhase}, nil
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

// --- Global mock instances ---

var (
	mockUsers    = &mockUserRepo{}
	mockPrograms = &mockProgramRepo{}
	mockPhases   = &mockPhaseRepo{}
)

func resetMocks() {
	*mockUsers = mockUserRepo{}
	*mockPrograms = mockProgramRepo{}
	*mockPhases = mockPhaseRepo{}
	lastProgramCreate = struct {
		name          string
		numPhases     int
		weeksPerPhase int
	}{}
	lastPhaseUpdates = nil
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
	name          string
	numPhases     int
	weeksPerPhase int
}

// captureProgramCreate makes CreateFn capture the call args and return a valid program.
func captureProgramCreate() {
	mockPrograms.CreateFn = func(userID int64, name string, startDate time.Time, numPhases, weeksPerPhase, defaultRepMin, defaultRepMax int) (*models.Program, error) {
		lastProgramCreate.name = name
		lastProgramCreate.numPhases = numPhases
		lastProgramCreate.weeksPerPhase = weeksPerPhase
		return &models.Program{ID: testProgramID, UserID: userID, Name: name, NumPhases: numPhases, WeeksPerPhase: weeksPerPhase}, nil
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
