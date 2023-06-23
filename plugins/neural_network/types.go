package neural_network

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
)

const (
	Name                = "neural_network"
	EntityNeuralNetwork = string("neural_network")
)

const (
	SettingParam1 = "param1"
	SettingParam2 = "param2"

	AttrPhase = "phase"

	StateEnabled  = "enabled"
	StateDisabled = "disabled"

	ActionEnabled = "enable"
	ActionDisable = "disable"

	Version = "0.0.1"
)

// store entity status in this struct
func NewAttr() m.Attributes {
	return m.Attributes{
		AttrPhase: {
			Name: AttrPhase,
			Type: common.AttributeString,
		},
	}
}

// entity settings
func NewSettings() m.Attributes {
	return m.Attributes{
		SettingParam1: {
			Name: SettingParam1,
			Type: common.AttributeString,
		},
		SettingParam2: {
			Name: SettingParam2,
			Type: common.AttributeString,
		},
	}
}

// state list entity
func NewStates() (states map[string]entity_manager.ActorState) {

	states = map[string]entity_manager.ActorState{
		StateEnabled: {
			Name:        StateEnabled,
			Description: "Enabled",
		},
		StateDisabled: {
			Name:        StateDisabled,
			Description: "disabled",
		},
	}

	return
}

// entity action list
func NewActions() map[string]entity_manager.ActorAction {
	return map[string]entity_manager.ActorAction{
		ActionEnabled: {
			Name:        ActionEnabled,
			Description: "enable",
		},
		ActionDisable: {
			Name:        ActionDisable,
			Description: "disable",
		},
	}
}
