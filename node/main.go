package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: ./node (Tracker IP)\n")
		os.Exit(1)
	}
	node := new(Node)

	node.Construct(os.Args[1])
	node.Start()
	/*time := time.Now().Format(time.RFC3339)
	 block := block.Block{
		Data:      "Hello World",
		Timestamp: time,
		UID:       "AZERTYUIOP"}

	block.BuildHash()
	block.SaveBlock("./") */
}
