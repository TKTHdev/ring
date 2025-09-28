package main

import (
	"fmt"
	"net/rpc"
	"slices"
	"sync"
	"time"
)

type Node struct {
	addr          string
	selfIdx       int
	nodeList      []string
	isLeader      bool
	leaderAddr    string
	rpcClients    map[string]*rpc.Client
	mu            sync.Mutex
	electionCh    chan ElectionArgs
	coordinatorCh chan CoordinatorArgs
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

func (n *Node) sendElectionToNode(addr string, aliveNodes []string, originaddr string) error {
	args := &ElectionArgs{AliveNodes: aliveNodes, OriginAddr: originaddr}
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

func (n *Node) sendElectionToNextAliveNode(aliveNodes []string, originAddr string) {
	nxtNodeIdx := (n.selfIdx + 1) % len(n.nodeList)
	for {
		if nxtNodeIdx == n.selfIdx {
			break
		}
		nxtNodeAddr := n.nodeList[nxtNodeIdx]
		err := n.sendElectionToNode(nxtNodeAddr, aliveNodes, originAddr)
		if err != nil {
			fmt.Println("Node", nxtNodeAddr, "is down. Trying next node...")
			nxtNodeIdx = (nxtNodeIdx + 1) % len(n.nodeList)
			continue
		}
		fmt.Println("Election message sent to", nxtNodeAddr)
		break
	}
}

func (n *Node) sendElectionToNextAliveNodeOrigin(aliveNodes []string) {
	nxtNodeIdx := (n.selfIdx + 1) % len(n.nodeList)
	for {
		if nxtNodeIdx == n.selfIdx {
			n.electionCh <- ElectionArgs{AliveNodes: aliveNodes, OriginAddr: n.addr}
		}
		nxtNodeAddr := n.nodeList[nxtNodeIdx]
		err := n.sendElectionToNode(nxtNodeAddr, aliveNodes, n.addr)
		if err != nil {
			fmt.Println("Node", nxtNodeAddr, "is down. Trying next node...")
			nxtNodeIdx = (nxtNodeIdx + 1) % len(n.nodeList)
			continue
		}
		fmt.Println("Election message sent to", nxtNodeAddr)
		break
	}
}

func (n *Node) startElection() {
	n.electionCh = make(chan ElectionArgs)
	go n.sendElectionToNextAliveNodeOrigin([]string{n.addr})
	select {
	case args := <-n.electionCh:
		fmt.Println("Election completed. Alive nodes:", args.AliveNodes)
		maxNode := slices.Max(args.AliveNodes)
		n.leaderAddr = maxNode
		n.isLeader = n.leaderAddr == n.addr
		return

	case <-time.After(2 * time.Second):
		n.startElection()
	}
}

func (n *Node) run() {
	go n.listenAndServeRPC()
	go n.connectToCluster()

	for {
		time.Sleep(500 * time.Millisecond)
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
