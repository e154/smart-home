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
	for _, raw := range []string{internetCheckAutomationRaw, hddCheckAutomationRaw} {
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
	internetCheckAutomationRaw = `{"id":"1","name":"internet_check","description":"every 5 sec","enabled":true,"condition":"and","triggers":[{"name":"every 1 minutes","plugin_name":"time","attributes":{"cron":{"name":"cron","type":1,"string":"*/5 * * * * *","array":[],"map":{}}},"entity":null,"script":null}],"conditions":[],"actions":[{"name":"ping","entity":null,"entity_id":"sensor.intermet_checker","entity_action_name":"PING","script":null},{"name":"web request","entity":null,"entity_id":"sensor.intermet_checker","entity_action_name":"CHECK","script":null}]}`
	hddCheckAutomationRaw = `{"id":"5","name":"hdd_check","description":"every 10 sec","enabled":true,"condition":"and","triggers":[{"name":"time","plugin_name":"time","attributes":{"cron":{"name":"cron","type":1,"string":"*/10 * * * * *","array":[],"map":{}}},"entity":null,"script":null}],"conditions":[],"actions":[{"name":"action","entity":null,"entity_id":"hdd.hdd1","entity_action_name":"CHECK","script":null}]}`
)
