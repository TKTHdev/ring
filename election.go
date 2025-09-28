package main

import (
	"fmt"
	"slices"
	"time"
)

func (n *Node) sendElectionToNode(addr string, aliveNodes []string, originAddr string) error {
	args := &ElectionArgs{AliveNodes: aliveNodes, OriginAddr: originAddr}
	reply := &ElectionReply{}
	err := n.sendRPC(addr, ElectionRPC, args, reply)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) startElection() {
	n.electionCh = make(chan ElectionArgs)
	go n.sendElectionToNextAliveNodeOrigin([]string{n.addr})
	select {
	case args := <-n.electionCh:
		fmt.Println("Election completed. Alive nodes:", args.AliveNodes)
		maxNode := slices.Max(args.AliveNodes)
		n.coordinatorAddr = maxNode
		n.isCoordinator = n.coordinatorAddr == n.addr
		return

	case <-time.After(2 * time.Second):
		n.startElection()
	}
}

func (n *Node) sendElectionToNextAliveNode(aliveNodes []string, originAddr string) {
	nxtNodeIdx := (n.selfIdx + 1) % len(n.nodeList)
	for {
		if n.nodeList[nxtNodeIdx] == n.nodeList[n.selfIdx] {
			return
		}
		nxtNodeAddr := n.nodeList[nxtNodeIdx]
		err := n.sendElectionToNode(nxtNodeAddr, aliveNodes, originAddr)
		if err != nil {
			fmt.Println("Node", nxtNodeAddr, "is down. Trying next node...")
			nxtNodeIdx = (nxtNodeIdx + 1) % len(n.nodeList)
			continue
		}
		fmt.Println("Election message sent to", nxtNodeAddr)
		return
	}
}

func (n *Node) sendElectionToNextAliveNodeOrigin(aliveNodes []string) {
	nxtNodeIdx := (n.selfIdx + 1) % len(n.nodeList)
	for {
		if n.nodeList[nxtNodeIdx] == n.nodeList[n.selfIdx] {
			n.electionCh <- ElectionArgs{AliveNodes: aliveNodes, OriginAddr: n.addr}
			return
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
