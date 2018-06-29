package main

import (
	"container/list"
	"fmt"
	"net"
	"os"
	"strings"
)

type Node struct {
	Listener net.Listener
	NodeList list.List
	Map      map[string]func(string)
}

func (node *Node) Construct() {
	listener, err := net.Listen("tcp", ":7179")
	node.Listener = listener
	if err != nil {
		fmt.Printf("Cannot listen on *:7179\n")
		os.Exit(1)
	}
	fmt.Printf("Listening on *:7179\n")
	node.NodeList.Init()
	node.Map["Node"] = node.addNodeToList
	node.Map["Block"] = node.addNodeToList
}

func (node *Node) addNodeToList(ip string) {
	node.NodeList.PushBack(ip)
}

func (node *Node) comunicate(conn net.Conn) {
	fmt.Printf("Connected to %s\n", conn.RemoteAddr().String())
	for {
		buf := make([]byte, 4096)
		len, err := conn.Read(buf)
		if err != nil || len == 0 {
			return
		}
		buffer := string(buf)
		key := strings.Split(buffer, " ")
		if _, ok := node.Map[key[0]]; ok {
			node.Map[key[0]](key[1])
		}
	}
}

func (node *Node) Start() {
	for {
		conn, err := node.Listener.Accept()
		if err != nil {
			fmt.Printf("Error on Accept\n")
			os.Exit(1)
		}
		go node.comunicate(conn)
	}
}
