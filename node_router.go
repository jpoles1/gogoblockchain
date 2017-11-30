package main

import (
	"net/http"
)

//Multinode routing
func registerNode(w http.ResponseWriter, r *http.Request) {
	//Parse/clean urls from GET?
	addr, ok := r.URL.Query()["addr"]
	if !ok || len(addr) < 1 {
		http.Error(w, "Err 400: URL Param 'addr' is missing", 400)
		return
	}
	uri, err := servBlockchain.registerNode(addr[0])
	if err != nil {
		w.Write([]byte("Error: " + err.Error()))
		return
	}
	w.Write([]byte("Added new blockchain node: " + uri))
}
func resolveNode(w http.ResponseWriter, r *http.Request) {
	msg := "Our chain is authoritative"
	if servBlockchain.resolveConflicts() {
		msg = "Our chain was replaced"
	}
	w.Write([]byte(msg))
}
