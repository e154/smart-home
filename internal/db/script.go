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
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/e154/smart-home/pkg/apperr"
	. "github.com/e154/smart-home/pkg/common"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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
					err = fmt.Errorf("%s: %w", fmt.Sprintf("script name \"%s\" not unique", script.Name), apperr.ErrScriptAdd)
					return
				}
			default:
				fmt.Printf("unknown code \"%s\"\n", pgErr.Code)
			}
		}
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrScriptAdd)
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
			err = fmt.Errorf("%s: %w", fmt.Sprintf("id \"%d\"", scriptId), apperr.ErrScriptNotFound)
			return
		}
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrScriptGet)
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
			err = fmt.Errorf("%s: %w", fmt.Sprintf("name \"%s\"", name), apperr.ErrScriptNotFound)
			return
		}
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrScriptGet)
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
					err = fmt.Errorf("%s: %w", fmt.Sprintf("script name \"%s\" not unique", script.Name), apperr.ErrScriptUpdate)
					return
				}
			default:
				fmt.Printf("unknown code \"%s\"\n", pgErr.Code)
			}
		}
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrScriptUpdate)
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
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrScriptVersionAdd)
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
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrScriptVersionDelete)
		return
	}

	return
}

// Delete ...
func (n Scripts) Delete(ctx context.Context, scriptId int64) (err error) {
	if err = n.DB(ctx).Delete(&Script{Id: scriptId}).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrScriptDelete)
	}
	return
}

// List ...
func (n *Scripts) List(ctx context.Context, limit, offset int, orderBy, sort string, query *string, ids *[]uint64) (list []*Script, total int64, err error) {

	list = make([]*Script, 0)
	q := n.DB(ctx).Model(Script{})
	if query != nil {
		q = q.Where("name ILIKE ? or source ILIKE ?", "%"+*query+"%", "%"+*query+"%")
	}
	if ids != nil {
		q = q.Where("id IN (?)", *ids)
	}
	if err = q.Count(&total).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrScriptList)
		return
	}
	err = q.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error
	if err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrScriptList)
	}
	return
}

// Search ...
func (n *Scripts) Search(ctx context.Context, query string, limit, offset int) (list []*Script, total int64, err error) {

	q := n.DB(ctx).Model(&Script{}).
		Where("name ILIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrScriptSearch)
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]*Script, 0)
	if err = q.Find(&list).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrScriptSearch)
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
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrScriptStat)
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
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrScriptStat)
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
