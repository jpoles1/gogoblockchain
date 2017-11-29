package main

//Node is a struct to represent a participant in the blockchain
type Node struct {
	addr string
}

//NodeSet is a struct to contain a set of unique Nodes
type NodeSet struct {
	set map[Node]bool
}

func (nodes *NodeSet) add(newNode Node) bool {
	_, found := nodes.set[newNode]
	nodes.set[newNode] = true
	return !found //False if it existed already
}
func (nodes *NodeSet) get(findNode Node) bool {
	_, found := nodes.set[findNode]
	return found //true if it existed already
}
