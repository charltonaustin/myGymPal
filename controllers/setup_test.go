package controllers_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/beego/beego/v2/server/web/session/postgres"
	"myGymPal/models"
	"myGymPal/routers"
)

var testProgramDate = time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)

func TestMain(m *testing.M) {
	// Use an absolute path so TestBeegoInit's internal Chdir + path join resolves correctly.
	root, _ := filepath.Abs("..")
	// Change to project root so template/migration paths resolve correctly.
	os.Chdir(root)
	// InitBeegoBeforeTest loads config, sets RunMode=test, and initialises templates + sessions.
	beego.InitBeegoBeforeTest(filepath.Join(root, "conf", "app.test.conf"))
	if err := models.Init(); err != nil {
		panic("failed to init DB: " + err.Error())
	}
	// Register routes AFTER config is loaded so routes inherit SessionOn=true.
	routers.Register()
	os.Exit(m.Run())
}
