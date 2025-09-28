package main

import (
	"fmt"
	"slices"
	"time"
)

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
