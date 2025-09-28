package main

import (
	"fmt"
	"net/rpc"
	"sync"
	"time"
)

type Node struct {
	addr            string
	selfIdx         int
	nodeList        []string
	isCoordinator   bool
	coordinatorAddr string
	rpcClients      map[string]*rpc.Client
	mu              sync.Mutex
	electionCh      chan ElectionArgs
	coordinatorCh   chan CoordinatorArgs
}

func NewNode() *Node {
	n := &Node{
		isCoordinator:   false,
		coordinatorAddr: "",
		mu:              sync.Mutex{},
		//those will be set later
		addr:       "",
		nodeList:   []string{},
		rpcClients: make(map[string]*rpc.Client),
	}
	n.readClusterConfigAndSet("cluster.conf")
	n.readNodeIndexAndSet()
	return n
}

func (n *Node) sendPingToNode(addr string) error {
	args := &PingArgs{}
	reply := &PingReply{}
	err := n.sendRPC(addr, PingRPC, args, reply)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) run() {
	go n.listenAndServeRPC()
	go n.connectToCluster()

	for {
		time.Sleep(100 * time.Millisecond)
		if n.isCoordinator {
			fmt.Println("I am the leader:", n.addr)
			continue
		}
		if n.coordinatorAddr == "" {
			n.startElection()
			n.announceCoordinator(n.coordinatorAddr)
		} else {
			err := n.sendPingToNode(n.coordinatorAddr)
			if err != nil {
				fmt.Println("Coordinator", n.coordinatorAddr, "is down. Starting election...")
				n.coordinatorAddr = ""
				n.startElection()
				fmt.Println("New coordinator elected:", n.coordinatorAddr)
				n.announceCoordinator(n.coordinatorAddr)
			} else {
				fmt.Println("Coordinator", n.coordinatorAddr, "is alive.")
			}
		}

	}
}
