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

func (n *Node) sendElectionToNode(addr string, aliveNodes []string, originAddr string) error {
	args := &ElectionArgs{AliveNodes: aliveNodes, OriginAddr: originAddr}
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
		time.Sleep(500 * time.Millisecond)
		if n.isCoordinator {
			fmt.Println("I am the leader:", n.addr)
			continue
		}
		if n.coordinatorAddr == "" {
			n.startElection()
		} else {
			err := n.sendPingToNode(n.coordinatorAddr)
			if err != nil {
				fmt.Println("Leader", n.coordinatorAddr, "is down. Starting election...")
				n.coordinatorAddr = ""
				n.startElection()
				n.announceCoordinator(n.coordinatorAddr)
			} else {
				fmt.Println("Leader", n.coordinatorAddr, "is alive.")
			}
		}

	}
}
