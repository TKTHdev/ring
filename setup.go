package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (n *Node) readClusterConfigAndSet(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			n.nodeList = append(n.nodeList, ":"+line)
		}
	}
}

func (n *Node) readNodeIndexAndSet() {
	if len(os.Args) < 1 {
		fmt.Println("Usage: go run main.go <node_index>")
		return
	}
	index, err := strconv.Atoi(os.Args[1])
	if err != nil || index < -1 || index >= len(n.nodeList) {
		fmt.Println(err)
		return
	}
	n.addr = n.nodeList[index]
}
