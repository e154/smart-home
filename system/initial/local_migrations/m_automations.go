package local_migrations

import (
	"context"
	"encoding/json"
	"github.com/e154/smart-home/api/dto"
	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/endpoint"

	"github.com/e154/smart-home/adaptors"
)

type MigrationAutomations struct {
	adaptors *adaptors.Adaptors
	endpoint *endpoint.Endpoint
}

func NewMigrationAutomations(adaptors *adaptors.Adaptors, endpoint *endpoint.Endpoint) *MigrationAutomations {
	return &MigrationAutomations{
		adaptors: adaptors,
		endpoint: endpoint,
	}
}

func (n *MigrationAutomations) Up(ctx context.Context, adaptors *adaptors.Adaptors) (err error) {
	if adaptors != nil {
		n.adaptors = adaptors
	}

	d := dto.NewDto()
	for _, raw := range []string{internetCheckAutomationRaw} {
		req := &api.NewTaskRequest{}
		if err = json.Unmarshal([]byte(raw), req); err != nil {
			return
		}
		task := d.Automation.ImportTask(req)
		if _, _, err = n.endpoint.Task.Import(ctx, task); err != nil {
			return err
		}
	}

	return
}

const (
	internetCheckAutomationRaw = `{"name":"task_internet_check","description":"","enabled":true,"condition":"and","triggers":[{"name":"trigger_internet_check","entity":{"id":"sensor.intermet_checker"},"script":{"id":15,"lang":"coffeescript","name":"internet_check","source":"\"use strict\";\n\nlastPing = null\n\n# entity\n# ##################################\ncheckStatus =(s)->   \n    res = http.get(s.host || \"https://google.com\")\n    if res.error \n        Actor.setState\n            new_state: 'NOT_CONNECTED'\n            storage_save: true\n        return\n\n    Actor.setState\n        new_state: 'CONNECTED'\n        storage_save: true\n\nping =(s)->\n  r = ExecuteSync 'ping', '-c', '1', s.ping_host || \"google.com\"\n  if r.err\n    #console.error \"Ошибка при выполнении команды ping: #{r.err}\"\n    return\n  pingIndex = r.out.indexOf('time=')\n\n  dir = 'up'\n  pingDifference = 0\n    \n  if pingIndex != -1\n    pingEndIndex = r.out.indexOf(' ms', pingIndex)\n    if pingEndIndex != -1\n      pingTime = parseFloat(r.out.substring(pingIndex + 5, pingEndIndex))\n      if lastPing\n        pingDifference = pingTime - lastPing\n        if lastPing > pingTime\n          dir = 'down'\n          #console.log \"Текущий пинг: #{pingTime} мс, разница с предыдущим пингом: #{pingDifference} мс\"\n      else\n        #console.log \"Текущий пинг: #{pingTime} мс\"\n      lastPing = pingTime\n  #else\n    #console.log 'Не удалось получить пинг'\n  \n  attrs =\n    ping: pingTime\n    dir: dir\n    diff: pingDifference\n  \n  # update storage if need\n  Actor.setState\n    attribute_values: attrs\n    storage_save: false\n\nentityAction = (entityId, actionName)->\n    entity = entityManager.getEntity(ENTITY_ID)\n    s = entity.getSettings()\n    switch actionName\n        when 'CHECK' then checkStatus(s)\n        when 'PING' then ping(s)\n\n# automation\n# ##################################\nautomationTriggerTime = (msg)->\n    entityManager.callAction(msg.entity_id, 'CHECK', {})\n    entityManager.callAction(msg.entity_id, 'PING', {})\n    return false\n","description":""},"plugin_name":"time","attributes":{"cron":{"name":"cron","type":1,"string":"*/5 * * * * *","array":[],"map":{}}}}],"conditions":[],"actions":[]}`
)
