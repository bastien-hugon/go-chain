package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

func getNodeList(ip string) ([]string, error) {
	ret := []string{}
	addr := strings.Join([]string{ip, "7180"}, ":")
	conn, err := net.Dial("tcp", addr)
	defer conn.Close()
	if err != nil {
		return ret, err
	}
	conn.Write([]byte("GetNodeList\n"))
	buff := make([]byte, 4096)
	n, err := conn.Read(buff)
	if err != nil {
		return ret, err
	}
	println("Received: " + string(buff[:n]))
	err = json.Unmarshal(buff[:n], &ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func main() {
	if len(os.Args) < 2 {
		os.Exit(84)
	}
	nodeList, err := getNodeList(os.Args[1])
	if err != nil {
		fmt.Errorf("Error while retrieving the node list")
		panic(err)
	}
	println("Node Ips: " + strings.Join(nodeList, ","))
}
