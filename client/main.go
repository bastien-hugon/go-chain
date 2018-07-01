package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func saveBlock(ip string, data string) error {
	addr := strings.Join([]string{ip, "7179"}, ":")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	fmt.Println("Send: data[" + data + "] to " + ip)
	defer conn.Close()
	_, err = conn.Write([]byte("Create " + data + "\n"))
	if err != nil {
		return err
	}
	return nil
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
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	base64Data := base64.StdEncoding.EncodeToString(bytes)
	fmt.Println("Data= " + base64Data)
	err = saveBlock(nodeList[0], base64Data)
	if err != nil {
		panic(err)
	}
}
