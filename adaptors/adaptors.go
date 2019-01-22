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
	Flow             *Flow
	FlowElement      *FlowElement
	Connection       *Connection
	Worker           *Worker
	Role             *Role
	Permission       *Permission
	User             *User
	UserMeta         *UserMeta
	Image            *Image
	Variable         *Variable
	Map              *Map
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
		Flow:             GetFlowAdaptor(db),
		FlowElement:      GetFlowElementAdaptor(db),
		Connection:       GetConnectionAdaptor(db),
		Worker:           GetWorkerAdaptor(db),
		Role:             GetRoleAdaptor(db),
		Permission:       GetPermissionAdaptor(db),
		User:             GetUserAdaptor(db),
		UserMeta:         GetUserMetaAdaptor(db),
		Image:            GetImageAdaptor(db),
		Variable:         GetVariableAdaptor(db),
		Map:              GetMapAdaptor(db),
	}

	return
}
