package initial

import (
	"github.com/e154/smart-home/system/migrations"
	"github.com/op/go-logging"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"fmt"
	"github.com/e154/smart-home/system/initial/env1"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/access_list"
)

var (
	log = logging.MustGetLogger("initial")
)

type InitialService struct {
	migrations    *migrations.Migrations
	adaptors      *adaptors.Adaptors
	core          *core.Core
	scriptService *scripts.ScriptService
	accessList    *access_list.AccessListService
}

func NewInitialService(migrations *migrations.Migrations,
	adaptors *adaptors.Adaptors,
	core *core.Core,
	scriptService *scripts.ScriptService,
	accessList *access_list.AccessListService) *InitialService {
	return &InitialService{
		migrations:    migrations,
		adaptors:      adaptors,
		core:          core,
		scriptService: scriptService,
		accessList:    accessList,
	}
}

func (n *InitialService) Reset() {

	log.Info("full reset")

	n.migrations.Purge()

	env1.Init(n.adaptors, n.accessList, n.scriptService)

	fmt.Println()
	log.Info("complete")
}
