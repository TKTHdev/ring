package main

import (
	"fmt"
	"net/rpc"
	"sync"
	"time"
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

func (n *Node) sendPingToNode(addr string) error {
	args := &PingArgs{}
	reply := &PingReply{}
	err := n.sendRPC(addr, PingRPC, args, reply)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) startElection() {
	fmt.Println("starting election")
}

func (n *Node) run() {
	go n.listenAndServeRPC()
	go n.connectToCluster()
	for {
		time.Sleep(1 * time.Second)
		if n.isLeader {
			fmt.Println("I am the leader:", n.addr)
			continue
		}
		if n.leaderAddr == "" {
			n.startElection()
		} else {
			err := n.sendPingToNode(n.leaderAddr)
			if err != nil {
				fmt.Println("Leader", n.leaderAddr, "is down. Starting election...")
				n.leaderAddr = ""
				n.startElection()
			} else {
				fmt.Println("Leader", n.leaderAddr, "is alive.")
			}
		}

	}
}
