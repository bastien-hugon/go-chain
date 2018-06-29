package main

func main() {
	node := new(Node)

	node.Construct()
	node.Start()
	/*time := time.Now().Format(time.RFC3339)
	 block := block.Block{
		Data:      "Hello World",
		Timestamp: time,
		UID:       "AZERTYUIOP"}

	block.BuildHash()
	block.SaveBlock("./") */
}
