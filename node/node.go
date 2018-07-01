package main

import (
	"encoding/json"
	"fmt"
	"github.com/bastien-hugon/go-chain/data_structures"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Node struct {
	Listener net.Listener
	NodeList []string
	Map      map[string]func(string, net.Conn)
	Blocks   []block.Block
	Conns    []net.Conn
}

func (node *Node) addNodeToList(data string, conn net.Conn) {
	ip := strings.Split(data, " ")[1]
	newIP := strings.Split(ip, ":")

	for i := 0; i < len(node.NodeList); i++ {
		if node.NodeList[i] == ip || node.NodeList[i] == newIP[0] {
			return
		}
	}
	if newIP[0] != "" {
		node.NodeList = append(node.NodeList, newIP[0])
	} else {
		node.NodeList = append(node.NodeList, ip)
	}
}

func (node *Node) broadcastBlock(data string) {
	for i := 0; i < len(node.Conns); i++ {
		node.Conns[i].Write([]byte(data))
	}
}

func (node *Node) saveNewBlock(data string, conn net.Conn) {
	var block block.Block

	line := []byte(data)[6:]
	err := json.Unmarshal(line, &block)
	if err != nil {
		return
	}
	block.BuildHash()
	for i := 0; i < len(node.Blocks); i++ {
		if node.Blocks[i].Hash == block.Hash {
			return
		}
	}
	block.SaveBlock("./blocks/")
	node.broadcastBlock(data)
	node.Blocks = append(node.Blocks, block)
}

func (node *Node) createBlock(data string, conn net.Conn) {
	block := new(block.Block)
	block.Data = strings.Split(data, " ")[1]
	block.UID = conn.RemoteAddr().String()
	block.Timestamp = strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	block.BuildHash()
	node.Blocks = append(node.Blocks, *block)
	err := block.SaveBlock("./blocks/")
	if err != nil {
		panic(err)
	}
	json, _ := json.Marshal(block)
	node.broadcastBlock(string(json))
	fmt.Printf("%s\n", string(json))
}

func (node *Node) refreshNodeList(conn net.Conn) {
	for {
		buf := make([]byte, 4096)
		len, err := conn.Read(buf)
		ret := []string{}
		err = json.Unmarshal(buf[:len], &ret)
		if err != nil {
			return
		}
		node.NodeList = ret
		fmt.Println("Node Ips: " + strings.Join(ret, ","))
		for _, v := range node.NodeList {
			conn, err := net.Dial("tcp", v+":7179")
			if err == nil {
				node.Conns = append(node.Conns, conn)
			}
		}
	}
}

func (node *Node) connectTracker(tracker string) {
	conn, err := net.Dial("tcp", tracker+":7180")
	if err == nil {
		fmt.Fprintf(conn, "Register\n")
		buf := make([]byte, 4096)
		len, err := conn.Read(buf)
		if err != nil || len == 0 {
			fmt.Printf("Tracker error on %s:7180\n", tracker)
			os.Exit(1)
		}
		fmt.Fprintf(conn, "GetNodeList\n")
		go node.refreshNodeList(conn)
	} else {
		fmt.Printf("No Tracker Found on %s:7180\n", tracker)
		os.Exit(1)
	}
}

func (node *Node) Construct(tracker string) {
	// Opening all block and load them in the stack
	files, err := ioutil.ReadDir("../blocks/")
	if err != nil {
		fmt.Printf("No block loaded\n")
	} else {
		for _, f := range files {
			tmpBlock, err := block.LoadBlock("../blocks/" + f.Name())
			if err == nil {
				node.Blocks = append(node.Blocks, tmpBlock)
			}
		}
	}

	// Opening the TCP server
	listener, err := net.Listen("tcp", ":7179")
	node.Listener = listener
	if err != nil {
		fmt.Printf("Cannot listen on *:7179\n")
		os.Exit(1)
	}
	fmt.Printf("Listening on *:7179\n")

	// Fill the command Map
	node.Map = make(map[string]func(string, net.Conn))
	node.Map["Block"] = node.saveNewBlock
	node.Map["Node"] = node.addNodeToList
	node.Map["Create"] = node.createBlock
	node.Map["OK"] = nil
	node.Map["KO"] = nil

	// Connect to the tracker
	go node.connectTracker(tracker)
}

func (node *Node) comunicate(conn net.Conn) {
	fmt.Printf("Connected to %s\n", conn.RemoteAddr().String())
	for {
		buf := make([]byte, 4096)
		len, err := conn.Read(buf)
		if err != nil || len == 0 {
			conn.Close()
			return
		}
		buf = buf[:len]
		buffer := string(buf)
		buffer = strings.Trim(buffer, "\n")
		key := strings.Split(buffer, " ")
		if _, ok := node.Map[key[0]]; ok {
			node.Map[key[0]](buffer, conn)
		} else {
			conn.Write([]byte("Bad command.\n"))
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
		node.Conns = append(node.Conns, conn)
		go node.comunicate(conn)
	}
}
