package models

import (
	"encoding/json"
	"time"
	"github.com/e154/smart-home/system/uuid"
)

type FlowElementModel struct {
	Uuid          uuid.UUID                 `json:"uuid"`
	Name          string                    `json:"name" valid:"MaxSize(254);Required"`
	Description   string                    `json:"description"`
	FlowId        int64                     `json:"flow_id" valid:"Required"`
	Script        *Script                   `json:"script"`
	ScriptId      *int64                    `json:"script_id"`
	Status        StatusType                `json:"status" valid:"Required"`
	FlowLink      *int64                    `json:"flow_link"`
	PrototypeType FlowElementsPrototypeType `json:"prototype_type" valid:"Required"`
	GraphSettings json.RawMessage           `json:"graph_settings"`
	CreatedAt     time.Time                 `json:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at"`
}