package main

import (
	"fmt"
	"time"
)

func main() {
	n := NewNode()
	go n.listenAndServeRPC()
	go n.connectToCluster()
	for {
		time.Sleep(1 * time.Second)
		fmt.Println(n.rpcClients)
	}
}
