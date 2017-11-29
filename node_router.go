package main

import "net/http"

//Multinode routing
func registerNode(w http.ResponseWriter, r *http.Request) {

}
func resolveNode(w http.ResponseWriter, r *http.Request) {
	servBlockchain.resolveConflicts()
	//servBlockchain.validChain(servBlockchain.Chain)
}
