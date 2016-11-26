package core

import (
	"encoding/json"
	"reflect"
	"../stream"
)

func GetNodesStatus() (result map[int64]string) {
	result = make(map[int64]string)
	for _, node := range bpmsPtr.nodes {
		result[node.Id] = node.GetConnectStatus()
	}

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