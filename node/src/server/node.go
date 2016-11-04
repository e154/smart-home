package server

type Node struct {}

func (n *Node) Version(request *string, version *string) error {

	*version = "0.1.0"

	return nil
}