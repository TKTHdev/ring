package main

import (
	"net"
	"net/rpc"
)

func (n *Node) listenAndServeRPC() {
	err := rpc.Register(n)
	if err != nil {
		return
	}
	l, err := net.Listen("tcp", n.addr)
	if err != nil {
		return
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func (n *Node) connectToCluster() {
	for {
		for _, addr := range n.nodeList {
			if addr != n.addr {
				n.connectToNode(addr)
			}
		}
	}
}

func (n *Node) connectToNode(addr string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if _, exists := n.rpcClients[addr]; !exists {
		client, err := rpc.Dial("tcp", addr)
		if err != nil {
			return
		}
		n.rpcClients[addr] = client
	}
}
