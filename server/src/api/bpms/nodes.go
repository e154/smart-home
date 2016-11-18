package bpms

import (
	"log"
	"time"
	"encoding/json"
	"reflect"
	"../models"
	"../stream"
)

func (b *BPMS) InitNodes() (err error) {
	var nodes []*models.Node
	b.nodes = make(map[int64]*models.Node)
	b.chanals = make(map[int64]chan string)

	log.Println("--------------------- NODES ---------------------")
	if nodes, err = models.GetAllEnabledNodes(); err != nil {
		return
	}

	for _, node := range nodes {
		b.AddNode(node)
	}

	if len(b.nodes) == 0 {
		return
	}

	return
}

func (b *BPMS) AddNode(node *models.Node) (err error) {
	if _, ok := b.nodes[node.Id]; ok {
		return
	}

	b.nodes[node.Id] = node
	b.chanals[node.Id] = make(chan string)

	go func(ch <- chan string) {
		var quit, disconnect bool
		connect := true
		for ;; {

			select {
			case c := <- ch:
				switch c {
				case "quit":
					quit = true
				case "disconnect":
					disconnect = true
				case "connect":
					connect = true
				default:

				}

			default:

			}

			if quit {
				node.TcpClose()
				break
			}

			if node.Errors > 5 {
				connect = true
			}

			if disconnect {
				node.TcpClose()
				connect = false
				disconnect = false
			}

			if connect {
				disconnect = false
				node.TcpClose()

				if _, err := node.RpcDial(); err == nil {
					node.Errors = 0
					//log.Printf("Node dial tcp %s:%d ... ok",node.Ip, node.Port)
					connect = false
					node.SetConnectStatus("connected")
				} else {
					//log.Printf("Node error %s", err.Error())
					node.SetConnectStatus("error")
				}
			}

			time.Sleep(time.Second)
		}
	}(b.chanals[node.Id])

	return
}

func (b *BPMS) RemoveNode(node *models.Node) (err error) {

	if _, ok := b.nodes[node.Id]; ok {
		b.chanals[node.Id] <- "quit"
		close(b.chanals[node.Id])
		delete(b.chanals, node.Id)
		delete(b.nodes, node.Id)
	}

	BroadcastNodesStatus()

	return
}

func (b *BPMS) ReloadNode(node *models.Node) (err error) {

	if _, ok := b.nodes[node.Id]; !ok {
		b.AddNode(node)
		return
	}

	b.nodes[node.Id].Status = node.Status
	b.nodes[node.Id].Ip = node.Ip
	b.nodes[node.Id].Port = node.Port
	b.nodes[node.Id].SetConnectStatus("wait")
	if node.Status == "disabled" {
		b.chanals[node.Id] <- "disconnect"
	} else {
		b.chanals[node.Id] <- "connect"
	}

	BroadcastNodesStatus()

	return
}

func (b *BPMS) ConnectNode(node *models.Node) (err error) {

	if _, ok := b.nodes[node.Id]; ok {
		b.chanals[node.Id] <- "connect"
	}

	BroadcastNodesStatus()

	return
}

func (b *BPMS) DisconnectNode(node *models.Node) (err error) {

	if _, ok := b.nodes[node.Id]; ok {
		b.chanals[node.Id] <- "disconnect"
	}

	BroadcastNodesStatus()

	return
}

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