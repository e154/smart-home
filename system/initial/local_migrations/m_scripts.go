package local_migrations

import (
	"context"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
	"github.com/e154/smart-home/system/scripts"
)

type MigrationScripts struct {
	adaptors *adaptors.Adaptors
	scriptService scripts.ScriptService
}

func NewMigrationScripts(adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService) *MigrationScripts {
	return &MigrationScripts{
		adaptors: adaptors,
		scriptService: scriptService,
	}
}

func (s *MigrationScripts) addScripts() (scripts []*m.Script, err error) {

	scripts = []*m.Script{}
	return
}

func (s *MigrationScripts) addScript(name, source, desc string) (script *m.Script, err error) {

	if script, err = s.adaptors.Script.GetByName(name); err == nil {
		return
	}

	script = &m.Script{
		Lang:        common.ScriptLangCoffee,
		Name:        name,
		Source:      source,
		Description: desc,
	}

	engineScript, err := s.scriptService.NewEngine(script)

	err = engineScript.Compile()
	So(err, ShouldBeNil)

	script.Id, err = s.adaptors.Script.Add(script)
	So(err, ShouldBeNil)
	return
}

func (n *MigrationScripts) Up(ctx context.Context, adaptors *adaptors.Adaptors) error {
	if adaptors != nil {
		n.adaptors = adaptors
	}
	_, err := n.addScripts()
	return err
}

