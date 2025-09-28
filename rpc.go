package main

import "fmt"

const (
	ElectionRPC    = "Node.Election"
	CoordinatorRPC = "Node.Coordinator"
	PingRPC        = "Node.Ping"
)

func (n *Node) Election(args *ElectionArgs, reply *struct{}) error {
	fmt.Println("Received election message from", args.OriginAddr)
	if args.OriginAddr == n.addr {
		n.electionCh <- *args
		return nil
	}
	args.AliveNodes = append(args.AliveNodes, n.addr)
	n.sendElectionToNextAliveNode(args.AliveNodes, args.OriginAddr)

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
