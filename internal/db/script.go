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
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/e154/smart-home/pkg/apperr"
	. "github.com/e154/smart-home/pkg/common"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Scripts ...
type Scripts struct {
	*Common
}

type ScriptInfo struct {
	AlexaIntents         int `gorm:"->"`
	EntityActions        int `gorm:"->"`
	EntityScripts        int `gorm:"->"`
	AutomationTriggers   int `gorm:"->"`
	AutomationConditions int `gorm:"->"`
	AutomationActions    int `gorm:"->"`
}

// Script ...
type Script struct {
	ScriptInfo
	Id          int64 `gorm:"primary_key"`
	Lang        ScriptLang
	Name        string
	Source      string
	Description string
	Compiled    string
	Versions    []*ScriptVersion
	CreatedAt   time.Time `gorm:"<-:create"`
	UpdatedAt   time.Time
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
	if err = n.DB(ctx).Create(&script).Error; err != nil {
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
	err = n.DB(ctx).Model(script).
		Where("id = ?", scriptId).
		Preload("Versions").
		First(&script).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrScriptNotFound, fmt.Sprintf("id \"%d\"", scriptId))
			return
		}
		err = errors.Wrap(apperr.ErrScriptGet, err.Error())
	}

	err = n.DB(ctx).Raw(`
	select
	      (select count(*) from alexa_intents where script_id = scripts.id)  as alexa_intents,
	      (select count(*) from entity_actions where script_id = scripts.id) as entity_actions,
	      (select count(*) from entity_scripts where script_id = scripts.id) as entity_scripts,
	      (select count(*) from triggers where script_id = scripts.id)       as automation_triggers,
	      (select count(*) from conditions where script_id = scripts.id)     as automation_conditions,
	      (select count(*) from actions where script_id = scripts.id)        as automation_actions
	from scripts where id = ?`, scriptId).
		First(&script.ScriptInfo).
		Error

	return
}

// GetByName ...
func (n Scripts) GetByName(ctx context.Context, name string) (script *Script, err error) {
	script = &Script{}
	err = n.DB(ctx).Model(script).
		Where("name = ?", name).
		Preload("Versions").
		First(&script).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrScriptNotFound, fmt.Sprintf("name \"%s\"", name))
			return
		}
		err = errors.Wrap(apperr.ErrScriptGet, err.Error())
	}

	err = n.DB(ctx).Raw(`
	select 
	      (select count(*) from alexa_intents where script_id = scripts.id)  as alexa_intents,
	      (select count(*) from entity_actions where script_id = scripts.id) as entity_actions,
	      (select count(*) from entity_scripts where script_id = scripts.id) as entity_scripts,
	      (select count(*) from triggers where script_id = scripts.id)       as automation_triggers,
	      (select count(*) from conditions where script_id = scripts.id)     as automation_conditions,
	      (select count(*) from actions where script_id = scripts.id)        as automation_actions
	from scripts where name = ?`, name).
		First(&script.ScriptInfo).
		Error

	return
}

// Update ...
func (n Scripts) Update(ctx context.Context, script *Script) (err error) {
	err = n.DB(ctx).Model(&Script{Id: script.Id}).Updates(map[string]interface{}{
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
		return
	}

	hash := md5.Sum([]byte(script.Source))
	version := &ScriptVersion{
		Lang:     script.Lang,
		Source:   script.Source,
		ScriptId: script.Id,
		Sum:      []byte(hex.EncodeToString(hash[:])),
	}
	if err = n.DB(ctx).Create(version).Error; err != nil {
		err = errors.Wrap(apperr.ErrScriptVersionAdd, err.Error())
		return
	}

	q := `delete from script_versions
where id not in (
    select id
    from script_versions
    where script_id = ?
    order by created_at desc
    limit 10
) and script_id = ?`
	if _, err = n.DB(ctx).Raw(q, script.Id, script.Id).Rows(); err != nil {
		err = errors.Wrap(apperr.ErrScriptVersionDelete, err.Error())
		return
	}

	return
}

// Delete ...
func (n Scripts) Delete(ctx context.Context, scriptId int64) (err error) {
	if err = n.DB(ctx).Delete(&Script{Id: scriptId}).Error; err != nil {
		err = errors.Wrap(apperr.ErrScriptDelete, err.Error())
	}
	return
}

// List ...
func (n *Scripts) List(ctx context.Context, limit, offset int, orderBy, sort string, query *string, ids *[]uint64) (list []*Script, total int64, err error) {

	list = make([]*Script, 0)
	q := n.DB(ctx).Model(Script{})
	if query != nil {
		q = q.Where("name LIKE ? or source LIKE ?", "%"+*query+"%", "%"+*query+"%")
	}
	if ids != nil {
		q = q.Where("id IN (?)", *ids)
	}
	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrScriptList, err.Error())
		return
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

	q := n.DB(ctx).Model(&Script{}).
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
	err = n.DB(ctx).Raw(`
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
	err = n.DB(ctx).Raw(`
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
