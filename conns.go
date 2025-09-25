package main

import (
	"fmt"
	"net"
	"net/rpc"
)

func (n *Node) listenAndServeRPC() {
	err := rpc.Register(n)
	if err != nil {
		fmt.Println("Error registering RPC:", err)
		return
	}
	l, err := net.Listen("tcp", n.addr)
	if err != nil {
		fmt.Println("Error starting RPC server:", err)
		return
	}
	defer l.Close()
	fmt.Println("RPC server listening on", n.addr)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
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
			//fmt.Println("Error connecting to node", addr, ":", err)
			return
		}
		n.rpcClients[addr] = client
		fmt.Println("Connected to node", addr)
	}
}
