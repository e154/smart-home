package core

func GetNode(id int64) *NodeBind {

	//TODO sync atomic
	nodes := corePtr.GetNodes()
	node, ok := nodes[id]
	if !ok {
		return nil
	}

	return &NodeBind{node:node}
}

func GetNodeList() (nodes []*NodeBind) {

	//TODO sync atomic
	nodes = []*NodeBind{}
	for _, node := range corePtr.GetNodes() {
		nodes = append(nodes, &NodeBind{node: node})
	}

	return
}

func GetFlow(id int64) interface{} {

	flow, err := corePtr.GetFlow(id)
	if err != nil {
		return nil
	} else {
		return flow
	}
}