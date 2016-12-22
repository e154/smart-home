package core

import (
	"reflect"
	"encoding/json"
	"../stream"
	"../models"
	"../log"
	"github.com/astaxie/beego/orm"
	"fmt"
)

// ------------------------------------------------
// Node
// ------------------------------------------------
func GetNodesStatus() (result map[int64]string) {
	result = make(map[int64]string)
	for _, node := range corePtr.nodes {
		result[node.Id] = node.GetConnectStatus()
	}

	return
}

func streamWorkflowsStatus(client *stream.Client, value interface{}) {

	return
}

func streamNodesStatus(client *stream.Client, value interface{}) {
	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	result := GetNodesStatus()
	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "nodes": result})
	client.Send(string(msg))
}

func BroadcastNodesStatus() {
	result := GetNodesStatus()
	msg, _ := json.Marshal(map[string]interface{}{"type": "broadcast", "value": &map[string]interface{}{"type": "nodes", "body": result}})
	Hub.Broadcast(string(msg))
}

// ------------------------------------------------
// Worker
// ------------------------------------------------
func streamDoWorker(client *stream.Client, value interface{}) {
	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	var worker_id float64
	var err error

	if worker_id, ok = v["worker_id"].(float64); !ok {
		log.Warn("bad id param")
		return
	}

	var worker *models.Worker
	if worker, err = models.GetWorkerById(int64(worker_id)); err != nil {
		client.Notify("error", err.Error())
		return
	}

	if err = corePtr.DoWorker(worker); err != nil {
		client.Notify("error", err.Error())
		return
	}

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "status": "ok"})
	client.Send(string(msg))
}

// ------------------------------------------------
// Action
// ------------------------------------------------
func streamDoAction(client *stream.Client, value interface{}) {
	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	var device_action_id, device_id float64
	var err error

	if device_action_id, ok = v["action_id"].(float64); !ok {
		log.Warn("bad device_action_id param")
		return
	}

	fmt.Println(reflect.TypeOf(v["device_id"]))
	if device_id, ok = v["device_id"].(float64); !ok {
		log.Warn("bad device_id param")
		return
	}

	var device_action *models.DeviceAction
	if device_action, err = models.GetDeviceActionById(int64(device_action_id)); err != nil {
		client.Notify("error", err.Error())
		return
	}

	var device *models.Device
	if device, err = models.GetDeviceById(int64(device_id)); err != nil {
		client.Notify("error", err.Error())
		return
	}

	// get device
	// ------------------------------------------------
	var devices []*models.Device
	if device.Address != nil {
		devices = append(devices, device)
	} else {
		// значит тут группа устройств
		var childs []*models.Device
		if childs, _, err = device_action.Device.GetChilds(); err != nil {
			return
		}

		for _, child := range childs {
			if child.Address == nil || child.Status != "enabled" {
				continue
			}

			device := &models.Device{}
			*device = *device_action.Device
			device.Id = child.Id
			device.Name = child.Name
			device.Address = new(int)
			*device.Address = *child.Address
			device.Device = &models.Device{Id:int64(device_id)}
			device.Tty = child.Tty
			devices = append(devices, device)
		}
	}

	// get node
	// ------------------------------------------------
	nodes := corePtr.GetNodes()
	var node *models.Node
	if _, ok := nodes[device_action.Device.Node.Id]; ok {
		node = nodes[device_action.Device.Node.Id]
	} else {
		// autoload nodes
		if node, err = models.GetNodeById(device_action.Device.Node.Id); err != nil {
			client.Notify("error", err.Error())
			return
		}

		if err = CorePtr().AddNode(node); err != nil {
			client.Notify("error", err.Error())
			return
		}
	}

	// get script
	// ------------------------------------------------
	o := orm.NewOrm()
	if _, err = o.LoadRelated(device_action, "Script"); err != nil {
		client.Notify("error", err.Error())
		return
	}

	for _, device := range devices {
		var action *Action
		if action, err = NewAction(device, device_action.Script, node); err != nil {
			log.Error(err.Error())
			client.Notify("error", err.Error())
			continue
		}

		body, _ := action.Do()
		client.Notify("success", body)
	}

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "status": "ok"})
	client.Send(string(msg))
}