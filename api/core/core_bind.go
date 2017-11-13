package core

func GetNode(id int64) interface{} {

	nodes := corePtr.GetNodes()
	node, ok := nodes[id]
	if ok {
		return node
	} else {
		return nil
	}

	return node
}

func GetFlow(id int64) interface{} {

	flow, err := corePtr.GetFlow(id)
	if err != nil {
		return nil
	} else {
		return flow
	}
}