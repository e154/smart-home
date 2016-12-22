package core

import (
	"reflect"
	"encoding/json"
	"../stream"
	"../models"
	"../log"
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
		msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "error": err.Error()})
		client.Send(string(msg))
		return
	}

	if err = corePtr.DoWorker(worker); err != nil {
		msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "error": err.Error()})
		client.Send(string(msg))
		return
	}

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "status": "ok"})
	client.Send(string(msg))
}