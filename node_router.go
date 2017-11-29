package main

import (
	"fmt"
	"net/http"
)

//Multinode routing
func registerNode(w http.ResponseWriter, r *http.Request) {
	addr, ok := r.URL.Query()["addr"]
	fmt.Println(addr)
	if !ok || len(addr) < 1 {
		http.Error(w, "Err 400: URL Param 'addr' is missing", 400)
		return
	}
	fmt.Println(servBlockchain.Nodes)
	servBlockchain.Nodes.add(Node{addr[0]})
	fmt.Println(servBlockchain.Nodes)
	w.Write([]byte("Added new blockchain node: " + addr[0]))
}
func resolveNode(w http.ResponseWriter, r *http.Request) {
	msg := "Our chain is authoritative"
	if servBlockchain.resolveConflicts() {
		msg = "Our chain was replaced"
	}
	w.Write([]byte(msg))
}
