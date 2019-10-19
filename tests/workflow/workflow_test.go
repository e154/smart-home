package workflow

import (
	"github.com/e154/smart-home/system/dig"
	l "github.com/e154/smart-home/system/logging"
	"github.com/e154/smart-home/system/migrations"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"testing"
	"time"
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
	err := container.Invoke(func(migrations *migrations.Migrations,
		lx *logrus.Logger,
		back *l.LogBackend) {

		l.Initialize(back)

		time.Sleep(time.Millisecond * 500)

		os.Exit(m.Run())
	})

	if err != nil {
		print(err.Error())
	}
}
