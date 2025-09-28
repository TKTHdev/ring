package main

import (
	"fmt"
	"time"
)

func (n *Node) sendCoordinatorToNode(addr, coordAddr, originAddr string) error {
	args := &CoordinatorArgs{CoordinatorAddr: coordAddr, OriginAddr: originAddr}
	reply := &CoordinatorReply{}
	err := n.sendRPC(addr, CoordinatorRPC, args, reply)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) announceCoordinator(addr string) {
	fmt.Println("Announcing new coordinator:", addr, "by node:", n.addr)
	n.coordinatorCh = make(chan CoordinatorArgs)
	go n.sendCoordinatorToNextAliveNodeOrigin(addr)
	select {
	case args := <-n.coordinatorCh:
		fmt.Println("Coordinator announcement completed. New coordinator:", args.CoordinatorAddr)
		return
	case <-time.After(2 * time.Second):
		return
	}
}

func (n *Node) sendCoordinatorToNextAliveNodeOrigin(coordAddr string) {
	nxtNodeIdx := (n.selfIdx + 1) % len(n.nodeList)
	for {
		if nxtNodeIdx == n.selfIdx {
			n.coordinatorCh <- CoordinatorArgs{CoordinatorAddr: coordAddr, OriginAddr: n.addr}
		}
		nxtNodeAddr := n.nodeList[nxtNodeIdx]
		err := n.sendCoordinatorToNode(nxtNodeAddr, coordAddr, n.addr)
		if err != nil {
			fmt.Println("Node", nxtNodeAddr, "is down. Trying next node...")
			nxtNodeIdx = (nxtNodeIdx + 1) % len(n.nodeList)
			continue
		}
		fmt.Println("Coordinator message sent to", nxtNodeAddr)
		break
	}
}

func (n *Node) sendCoordinatorToNextAliveNode(coordAddr, originAddr string) {
	nxtNodeIdx := (n.selfIdx + 1) % len(n.nodeList)
	for {
		if nxtNodeIdx == n.selfIdx {
			return
		}
		nxtNodeAddr := n.nodeList[nxtNodeIdx]
		err := n.sendCoordinatorToNode(nxtNodeAddr, coordAddr, originAddr)
		if err != nil {
			fmt.Println("Node", nxtNodeAddr, "is down. Trying next node...")
			nxtNodeIdx = (nxtNodeIdx + 1) % len(n.nodeList)
			continue
		}
		fmt.Println("Coordinator message sent to", nxtNodeAddr)
		break
	}
}
