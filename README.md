# Geode
A simple Blockchain in Golang
This Blockchain is constitued by 3 Parts: `Node`, `Tracker` and `Client`

## Node
> A node is a TCP Server that receive, send and create Blocks from the Blockchain
### Commands
```go
/
| // --> Send Command
| // <-- Receive Command
|
| // Register
| --> Register
| // Ask to have some Node IP
| --> GetNodeList
| // Receive a new Node IP the this Node.
| <-- Node `IP`
| --> Node `IP`
| // Sent by the client to create a new block
| <-- Create `BLOCK DATA (BASE64)`
| // Send / Receive a Block to this Node in Json
| --> Block `BLOCK DATA (JSON)`
| <-- Block `BLOCK DATA (JSON)`
\
```

## Tracker
> A Tracker is a TCP Server that receive, send and create Blocks from the Blockchain
> It also contain all Node IP
### Commands
```go
/
| // --> Send Command
| // <-- Receive Command
|
| // Add the sending node to the list.
| <-- Register
| // Send all Nodes IP
| <-- GetNodeList
| // Receive a new Node IP the this Node.
| <-- Node `IP`
| --> Node `IP`
| // Send / Receive a Block to this Node in Json
| --> Block `BLOCK DATA (JSON)`
| <-- Block `BLOCK DATA (JSON)`
\
```