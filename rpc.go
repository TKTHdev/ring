package main

import "fmt"

const (
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

func (n *Node) sendRPC(targetAddr string, method string, args interface{}, reply interface{}) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	client, ok := n.rpcClients[targetAddr]
	if !ok {
		return fmt.Errorf("no RPC client for address: %s", targetAddr)
	}
	if err := client.Call(method, args, reply); err != nil {
		n.rpcClients[targetAddr].Close()
		delete(n.rpcClients, targetAddr)
		return err
	}
	return nil
}
