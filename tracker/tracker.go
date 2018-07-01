package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type Tracker struct {
	Listener net.Listener
	NodeList []string
	Map      map[string]func(string, net.Conn)
}

func (tracker *Tracker) removeNodeFromList(cmd string, conn net.Conn) {
	ip := conn.RemoteAddr().String()
	ip = strings.Split(ip, ":")[0]
	for i, v := range tracker.NodeList {
		if v == ip {
			tracker.NodeList = append(tracker.NodeList[:i], tracker.NodeList[i+1:]...)
			break
		}
	}
	fmt.Printf("UnRegistered IP: %s\n", ip)
	conn.Write([]byte("OK\n"))
}

func (tracker *Tracker) addNodeToList(cmd string, conn net.Conn) {
	ip := conn.RemoteAddr().String()
	ip = strings.Split(ip, ":")[0]
	tracker.NodeList = append(tracker.NodeList, ip)
	fmt.Printf("Registered IP: %s\n", ip)
	conn.Write([]byte("OK\n"))
}

func (tracker *Tracker) sendNodeList(cmd string, conn net.Conn) {
	json, err := json.Marshal(tracker.NodeList)
	if err != nil {
		fmt.Printf("Error while sending the nodeList\n")
		conn.Write([]byte("KO\n"))
	}
	conn.Write(append(json, []byte("\n")...))
}

func (tracker *Tracker) construct() {
	listener, err := net.Listen("tcp", ":7180")
	tracker.Listener = listener
	if err != nil {
		fmt.Printf("Cannot listen on *:7180\n")
		os.Exit(1)
	}
	fmt.Printf("Listening on *:7180\n")
	tracker.Map = make(map[string]func(string, net.Conn))
	tracker.Map["Register"] = tracker.addNodeToList
	tracker.Map["GetNodeList"] = tracker.sendNodeList
	tracker.Map["Unregister"] = tracker.removeNodeFromList
}

func (tracker *Tracker) comunicate(conn net.Conn) {
	fmt.Printf("Connected to %s\n", conn.RemoteAddr().String())
	for {
		buf := make([]byte, 4096)
		len, err := conn.Read(buf)
		if err != nil || len == 0 {
			return
		}
		buf[len-1] = '\000'
		buffer := string(buf)
		buffer = strings.Trim(buffer, "\n")
		key := strings.Split(buffer, "\000")
		if _, ok := tracker.Map[key[0]]; ok {
			tracker.Map[key[0]](buffer, conn)
		} else {
			conn.Write([]byte("Bad command.\n"))
		}
	}
}

func (tracker *Tracker) start() {
	for {
		conn, err := tracker.Listener.Accept()
		if err != nil {
			fmt.Printf("Error on Accept\n")
			os.Exit(1)
		}
		go tracker.comunicate(conn)
	}
}
