package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

//Multinode routing
func registerNode(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params)
	addr, ok := r.URL.Query()["addr"]
	if !ok || len(addr) < 1 {
		http.Error(w, "Err 400: URL Param 'addr' is missing", 400)
		return
	}
	servBlockchain.Nodes.add(Node{addr[0]})
	w.Write([]byte("Added new blockchain node: " + addr[0]))
}
func resolveNode(w http.ResponseWriter, r *http.Request) {
	msg := "Our chain is authoritative"
	if servBlockchain.resolveConflicts() {
		msg = "Our chain was replaced"
	}
	w.Write([]byte(msg))
}
