package main

import (
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
	n.coordinatorCh = make(chan CoordinatorArgs)
	go n.sendCoordinatorToNextAliveNodeOrigin(addr)
	select {
	case <-n.coordinatorCh:
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
			nxtNodeIdx = (nxtNodeIdx + 1) % len(n.nodeList)
			continue
		}
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
			nxtNodeIdx = (nxtNodeIdx + 1) % len(n.nodeList)
			continue
		}
		break
	}
}
