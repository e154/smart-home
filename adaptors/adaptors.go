// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package adaptors

import (
	"context"

	"github.com/e154/smart-home/common/logger"

	"github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/orm"
	"github.com/jinzhu/gorm"
	"go.uber.org/fx"
)

var (
	log = logger.MustGetLogger("adaptors")
)

// Adaptors ...
type Adaptors struct {
	db                *gorm.DB
	isTx              bool
	Script            IScript
	Role              IRole
	Permission        IPermission
	User              IUser
	UserMeta          IUserMeta
	Image             IImage
	Variable          IVariable
	Map               IMap
	MapLayer          IMapLayer
	MapText           IMapText
	MapImage          IMapImage
	MapElement        IMapElement
	Entity            IEntity
	EntityState       IEntityState
	EntityAction      IEntityAction
	EntityStorage     IEntityStorage
	Log               ILog
	Template          ITemplate
	Message           IMessage
	MessageDelivery   IMessageDelivery
	Zigbee2mqtt       IZigbee2mqtt
	Zigbee2mqttDevice IZigbee2mqttDevice
	AlexaSkill        IAlexaSkill
	AlexaIntent       IAlexaIntent
	Storage           IStorage
	Metric            IMetric
	MetricBucket      IMetricBucket
	Area              IArea
	Action            IAction
	Condition         ICondition
	Trigger           ITrigger
	Task              ITask
	RunHistory        IRunHistory
	Plugin            IPlugin
	TelegramChat      ITelegramChat
	Dashboard         IDashboard
	DashboardTab      IDashboardTab
	DashboardCard     IDashboardCard
	DashboardCardItem IDashboardCardItem
}

// NewAdaptors ...
func NewAdaptors(lc fx.Lifecycle,
	db *gorm.DB,
	cfg *models.AppConfig,
	migrations *migrations.Migrations,
	orm *orm.Orm) (adaptors *Adaptors) {

	adaptors = &Adaptors{
		db:                db,
		Script:            GetScriptAdaptor(db),
		Role:              GetRoleAdaptor(db),
		Permission:        GetPermissionAdaptor(db),
		User:              GetUserAdaptor(db),
		UserMeta:          GetUserMetaAdaptor(db),
		Image:             GetImageAdaptor(db),
		Variable:          GetVariableAdaptor(db),
		Map:               GetMapAdaptor(db),
		MapLayer:          GetMapLayerAdaptor(db),
		MapText:           GetMapTextAdaptor(db),
		MapImage:          GetMapImageAdaptor(db),
		Entity:            GetEntityAdaptor(db),
		EntityState:       GetEntityStateAdaptor(db),
		EntityAction:      GetEntityActionAdaptor(db),
		EntityStorage:     GetEntityStorageAdaptor(db),
		MapElement:        GetMapElementAdaptor(db),
		Log:               GetLogAdaptor(db),
		Template:          GetTemplateAdaptor(db),
		Message:           GetMessageAdaptor(db),
		MessageDelivery:   GetMessageDeliveryAdaptor(db),
		Zigbee2mqtt:       GetZigbee2mqttAdaptor(db),
		Zigbee2mqttDevice: GetZigbee2mqttDeviceAdaptor(db),
		AlexaSkill:        GetAlexaSkillAdaptor(db),
		AlexaIntent:       GetAlexaIntentAdaptor(db),
		Storage:           GetStorageAdaptor(db),
		Metric:            GetMetricAdaptor(db, orm),
		MetricBucket:      GetMetricBucketAdaptor(db, orm),
		Area:              GetAreaAdaptor(db),
		Action:            GetActionAdaptor(db),
		Condition:         GetConditionAdaptor(db),
		Trigger:           GetTriggerAdaptor(db),
		Task:              GetTaskAdaptor(db),
		RunHistory:        GetRunHistoryAdaptor(db),
		Plugin:            GetPluginAdaptor(db),
		TelegramChat:      GetTelegramChannelAdaptor(db),
		Dashboard:         GetDashboardAdaptor(db),
		DashboardTab:      GetDashboardTabAdaptor(db),
		DashboardCard:     GetDashboardCardAdaptor(db),
		DashboardCardItem: GetDashboardCardItemAdaptor(db),
	}

	if lc != nil {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) (err error) {
				if cfg != nil && migrations != nil && cfg.AutoMigrate {
					err = migrations.Up()
					return
				}
				return
			},
		})
	}

	return
}

// Begin ...
func (a Adaptors) Begin() (adaptors *Adaptors) {
	adaptors = NewAdaptors(nil, a.db.Begin(), nil, nil, nil)
	adaptors.isTx = true
	return
}

// Commit ...
func (a *Adaptors) Commit() error {
	if !a.isTx {
		return nil
	}
	a.isTx = false
	return a.db.Commit().Error
}

// Rollback ...
func (a *Adaptors) Rollback() error {
	if !a.isTx {
		return nil
	}
	return a.db.Rollback().Error
}
