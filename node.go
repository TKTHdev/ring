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
