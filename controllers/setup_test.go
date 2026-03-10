package controllers_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"myGymPal/controllers"
	"myGymPal/routers"
)

var testProgramDate = time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)

func TestMain(m *testing.M) {
	root, _ := filepath.Abs("..")
	os.Chdir(root)
	beego.InitBeegoBeforeTest(filepath.Join(root, "conf", "app.test.conf"))
	// Register routes AFTER config is loaded so routes inherit SessionOn=true.
	routers.Register()
	// Override the real repositories with mocks.
	controllers.Users = mockUsers
	controllers.Programs = mockPrograms
	controllers.Phases = mockPhases
	os.Exit(m.Run())
}
