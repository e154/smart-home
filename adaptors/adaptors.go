package adaptors

import (
	"github.com/op/go-logging"
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/migrations"
)

var (
	log = logging.MustGetLogger("adaptors")
)

type Adaptors struct {
	Node             *Node
	Script           *Script
	Workflow         *Workflow
	WorkflowScenario *WorkflowScenario
	Device           *Device
	DeviceAction     *DeviceAction
	DeviceState      *DeviceState
}

func NewAdaptors(db *gorm.DB,
	cfg *config.AppConfig,
	migrations *migrations.Migrations) (adaptors *Adaptors) {

	if cfg.AutoMigrate {
		if err := migrations.Up(); err != nil {
			panic(err.Error())
		}
	}

	adaptors = &Adaptors{
		Node:             GetNodeAdaptor(db),
		Script:           GetScriptAdaptor(db),
		Workflow:         GetWorkflowAdaptor(db),
		WorkflowScenario: GetWorkflowScenarioAdaptor(db),
		Device:           GetDeviceAdaptor(db),
		DeviceAction:     GetDeviceActionAdaptor(db),
		DeviceState:      GetDeviceStateAdaptor(db),
	}

	return
}
