package main

import (
	"net/rpc"
	"sync"
)

type Node struct {
	addr       string
	nodeList   []string
	isLeader   bool
	leaderAddr string
	rpcClients map[string]*rpc.Client
	mu         sync.Mutex
}

func NewNode() *Node {
	n := &Node{
		isLeader:   false,
		leaderAddr: "",
		mu:         sync.Mutex{},
		//those will be set later
		addr:       "",
		nodeList:   []string{},
		rpcClients: make(map[string]*rpc.Client),
	}
	n.readClusterConfigAndSet("cluster.conf")
	n.readNodeIndexAndSet()
	return n
}

func (n *Node) sendElectionToNode(addr string, aliveNodes []string) error {
	args := &ElectionArgs{AliveNodes: aliveNodes}
	reply := &ElectionReply{}
	err := n.sendRPC(addr, ElectionRPC, args, reply)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) sendCoordinatorToNode(addr, coordAddr string) error {
	args := &CoordinatorArgs{CoordinatorAddr: coordAddr}
	reply := &CoordinatorReply{}
	err := n.sendRPC(addr, CoordinatorRPC, args, reply)
	if err != nil {
		return err
	}
	return nil
}
