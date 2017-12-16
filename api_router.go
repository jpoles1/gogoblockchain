package main

import (
	"fmt"
	"net/http"
)

//Main blockchain routing
func voteAPI(w http.ResponseWriter, r *http.Request) {
	pollid, ok := r.URL.Query()["pollid"]
	if !ok || len(pollid) < 1 {
		http.Error(w, "Err 400: URL Param 'pollid' is missing", 400)
		return
	}
	voterid, ok := r.URL.Query()["voterid"]
	if !ok || len(voterid) < 1 {
		http.Error(w, "Err 400: URL Param 'voterid' is missing", 400)
		return
	}
	vote, ok := r.URL.Query()["vote"]
	if !ok || len(vote) < 1 {
		http.Error(w, "Err 400: URL Param 'vote' is missing", 400)
		return
	}
	servBlockchain.newTransaction(Transaction{pollid[0], voterid[0], vote[0]})
	w.Write([]byte("Vote sucessfully recorded."))
	go servBlockchain.tryToMine()
}
func newPollAPI(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()                     // Parses the request body
	x := r.Form.Get("parameter_name") // x will be "" if parameter is not set
	fmt.Println(x)
	w.Write([]byte("Vote sucessfully recorded."))
}
