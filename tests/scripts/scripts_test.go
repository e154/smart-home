package scripts

import (
	"time"
	"os"
	"testing"
	"github.com/e154/smart-home/system/dig"
	"github.com/e154/smart-home/system/migrations"
	"path/filepath"
	"github.com/sirupsen/logrus"
	l "github.com/e154/smart-home/system/logging"
)

func init() {
	apppath := filepath.Join(os.Getenv("PWD"), "../..")
	os.Chdir(apppath)
}

var (
	container *dig.Container
)

func TestMain(m *testing.M) {

	container = BuildContainer()
	container.Invoke(func(migrations *migrations.Migrations,
		lx *logrus.Logger) {

		l.Initialize(lx)

		time.Sleep(time.Millisecond * 500)

		os.Exit(m.Run())
	})
}
