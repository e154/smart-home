package use_case

import (
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/adaptors"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("use_case")
)

type Command struct {
	Auth             *AuthCommand
	Device           *DeviceCommand
	DeviceAction     *DeviceActionCommand
	DeviceState      *DeviceStateCommand
	Flow             *FlowCommand
	Image            *ImageCommand
	Log              *LogCommand
	Map              *MapCommand
	MapElement       *MapElementCommand
	MapLayer         *MapLayerCommand
	Node             *NodeCommand
	Role             *RoleCommand
	Script           *ScriptCommand
	Workflow         *WorkflowCommand
	WorkflowScenario *WorkflowScenarioCommand
	User             *UserCommand
}

func NewUseCase(adaptors *adaptors.Adaptors,
	core *core.Core,
	scriptService *scripts.ScriptService,
	accessList *access_list.AccessListService) *Command {
	common := NewCommonCommand(adaptors, core, accessList, scriptService)
	return &Command{
		Auth:             NewAuthCommand(common),
		Device:           NewDeviceCommand(common),
		DeviceAction:     NewDeviceActionCommand(common),
		DeviceState:      NewDeviceStateCommand(common),
		Flow:             NewFlowCommand(common),
		Image:            NewImageCommand(common),
		Log:              NewLogCommand(common),
		Map:              NewMapCommand(common),
		MapElement:       NewMapElementCommand(common),
		MapLayer:         NewMapLayerCommand(common),
		Node:             NewNodeCommand(common),
		Role:             NewRoleCommand(common),
		Script:           NewScriptCommand(common),
		Workflow:         NewWorkflowCommand(common),
		WorkflowScenario: NewWorkflowScenarioCommand(common),
		User:             NewUserCommand(common),
	}
}
