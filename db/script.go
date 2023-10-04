// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
)

// Scripts ...
type Scripts struct {
	Db *gorm.DB
}

// Script ...
type Script struct {
	Id                   int64 `gorm:"primary_key"`
	Lang                 ScriptLang
	Name                 string
	Source               string
	Description          string
	Compiled             string
	AlexaIntents         int `gorm:"->"`
	EntityActions        int `gorm:"->"`
	EntityScripts        int `gorm:"->"`
	AutomationTriggers   int `gorm:"->"`
	AutomationConditions int `gorm:"->"`
	AutomationActions    int `gorm:"->"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type ScriptsStatistic struct {
	Total        int32
	Used         int32
	Unused       int32
	CoffeeScript int32
	TypeScript   int32
	JavaScript   int32
}

// TableName ...
func (d *Script) TableName() string {
	return "scripts"
}

// Add ...
func (n Scripts) Add(ctx context.Context, script *Script) (id int64, err error) {
	if err = n.Db.WithContext(ctx).Create(&script).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				if strings.Contains(pgErr.Message, "name_at_scripts_unq") {
					err = errors.Wrap(apperr.ErrScriptAdd, fmt.Sprintf("script name \"%s\" not unique", script.Name))
					return
				}
			default:
				fmt.Printf("unknown code \"%s\"\n", pgErr.Code)
			}
		}
		err = errors.Wrap(apperr.ErrScriptAdd, err.Error())
		return
	}
	id = script.Id
	return
}

// GetById ...
func (n Scripts) GetById(ctx context.Context, scriptId int64) (script *Script, err error) {
	script = &Script{}
	err = n.Db.WithContext(ctx).Raw(`
select scripts.*,
       (select count(*) from alexa_intents where script_id = scripts.id)  as alexa_intents,
       (select count(*) from entity_actions where script_id = scripts.id) as entity_actions,
       (select count(*) from entity_scripts where script_id = scripts.id) as entity_scripts,
       (select count(*) from triggers where script_id = scripts.id)       as automation_triggers,
       (select count(*) from conditions where script_id = scripts.id)     as automation_conditions,
       (select count(*) from actions where script_id = scripts.id)        as automation_actions
from scripts where id = ?`, scriptId).
		First(script).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrScriptNotFound, fmt.Sprintf("id \"%d\"", scriptId))
			return
		}
		err = errors.Wrap(apperr.ErrScriptGet, err.Error())
	}
	return
}

// GetByName ...
func (n Scripts) GetByName(ctx context.Context, name string) (script *Script, err error) {
	script = &Script{}
	err = n.Db.WithContext(ctx).Raw(`
select scripts.*,
       (select count(*) from alexa_intents where script_id = scripts.id)  as alexa_intents,
       (select count(*) from entity_actions where script_id = scripts.id) as entity_actions,
       (select count(*) from entity_scripts where script_id = scripts.id) as entity_scripts,
       (select count(*) from triggers where script_id = scripts.id)       as automation_triggers,
       (select count(*) from conditions where script_id = scripts.id)     as automation_conditions,
       (select count(*) from actions where script_id = scripts.id)        as automation_actions
from scripts where name = ?`, name).
		First(script).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrScriptNotFound, fmt.Sprintf("name \"%d\"", name))
			return
		}
		err = errors.Wrap(apperr.ErrScriptGet, err.Error())
	}
	return
}

// Update ...
func (n Scripts) Update(ctx context.Context, script *Script) (err error) {
	err = n.Db.WithContext(ctx).Model(&Script{Id: script.Id}).Updates(map[string]interface{}{
		"name":        script.Name,
		"description": script.Description,
		"lang":        script.Lang,
		"source":      script.Source,
		"compiled":    script.Compiled,
	}).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				if strings.Contains(pgErr.Message, "name_at_scripts_unq") {
					err = errors.Wrap(apperr.ErrScriptUpdate, fmt.Sprintf("script name \"%s\" not unique", script.Name))
					return
				}
			default:
				fmt.Printf("unknown code \"%s\"\n", pgErr.Code)
			}
		}
		err = errors.Wrap(apperr.ErrScriptUpdate, err.Error())
	}
	return
}

// Delete ...
func (n Scripts) Delete(ctx context.Context, scriptId int64) (err error) {
	if err = n.Db.WithContext(ctx).Delete(&Script{Id: scriptId}).Error; err != nil {
		err = errors.Wrap(apperr.ErrScriptDelete, err.Error())
	}
	return
}

// List ...
func (n *Scripts) List(ctx context.Context, limit, offset int, orderBy, sort string, query *string) (list []*Script, total int64, err error) {

	if err = n.Db.WithContext(ctx).Model(Script{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrScriptList, err.Error())
		return
	}

	list = make([]*Script, 0)
	q := n.Db
	if query != nil {
		q = q.Where("name LIKE ?", "%"+*query+"%")
	}

	err = q.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error
	if err != nil {
		err = errors.Wrap(apperr.ErrScriptList, err.Error())
	}
	return
}

// Search ...
func (n *Scripts) Search(ctx context.Context, query string, limit, offset int) (list []*Script, total int64, err error) {

	q := n.Db.WithContext(ctx).Model(&Script{}).
		Where("name LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrScriptSearch, err.Error())
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]*Script, 0)
	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrScriptSearch, err.Error())
	}
	return
}

// Statistic ...
func (n *Scripts) Statistic(ctx context.Context) (statistic *ScriptsStatistic, err error) {

	statistic = &ScriptsStatistic{}

	var usedList []struct {
		Count int32
		Used  bool
	}
	err = n.Db.WithContext(ctx).Raw(`
select count(scripts.id),
       (exists(select * from alexa_intents where script_id = scripts.id) or 
        exists(select * from entity_actions where script_id = scripts.id) or
        exists(select * from entity_scripts where script_id = scripts.id) or
        exists(select * from triggers where script_id = scripts.id)       or
        exists(select * from conditions where script_id = scripts.id)     or
        exists(select * from actions where script_id = scripts.id)    ) as used
from scripts
group by used`).
		Scan(&usedList).
		Error

	if err != nil {
		err = errors.Wrap(apperr.ErrScriptStat, err.Error())
		return
	}

	for _, item := range usedList {
		statistic.Total += item.Count
		if item.Used {
			statistic.Used = item.Count

			continue
		}
		statistic.Unused = item.Count
	}

	var langList []struct {
		Lang  string
		Count int32
	}
	err = n.Db.WithContext(ctx).Raw(`
select scripts.lang, count(scripts.*)
		from scripts
		group by lang`).
		Scan(&langList).
		Error

	if err != nil {
		err = errors.Wrap(apperr.ErrScriptStat, err.Error())
		return
	}

	for _, item := range langList {
		switch item.Lang {
		case "coffeescript":
			statistic.CoffeeScript = item.Count
		case "ts":
			statistic.TypeScript = item.Count
		case "javascript":
			statistic.JavaScript = item.Count
		}
	}

	return
}
