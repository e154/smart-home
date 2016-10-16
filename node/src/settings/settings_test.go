package settings_test

import (
	"testing"
	"fmt"
	"../settings"
)

func TestSettingsPtr(t *testing.T) {

	s := settings.SettingsPtr()
	if s == nil {
		t.Errorf("Settings pointer is nil")
	}
}

func TestSettings_AppVresion(t *testing.T) {

	s := settings.SettingsPtr()
	ver := s.AppVresion()
	cur_ver := fmt.Sprintf("%d.%d.%d", settings.APP_MAJOR, settings.APP_MINOR, settings.APP_PATCH)
	if ver != cur_ver {
		t.Errorf("Bad version %s != %s", ver, cur_ver)
	}
}

func TestSettings_Init(t *testing.T) {

	s := settings.SettingsPtr()
	if (s.Init() != s) {
		t.Errorf("Settings pointer is nil")
	}
}

func TestSettings_Load(t *testing.T) {

	s := settings.SettingsPtr()

	ns, err := s.Load()
	if err != nil {
		t.Errorf("error %s", err.Error())
	}

	if (ns != s) {
		t.Errorf("Settings pointer is nil")
	}
}

func TestSettings_Save(t *testing.T) {

	s := settings.SettingsPtr()

	sn, err := s.Save()
	if err != nil {
		t.Errorf("error %s", err.Error())
	}

	if (sn != s) {
		t.Errorf("Settings pointer is nil")
	}
}