package main

const (
	// RPC names
	ElectionRPC    = "Node.Election"
	CoordinatorRPC = "Node.Coordinator"
	PingRPC        = "Node.Ping"
)

func (n *Node) Election(args *struct{}, reply *struct{}) error {
	return nil
}

func (n *Node) Coordinator(args *struct{}, reply *struct{}) error {
	return nil
}

func (n *Node) Ping(args *struct{}, reply *struct{}) error {
	return nil
}
