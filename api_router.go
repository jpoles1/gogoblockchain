package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

//Main blockchain routing
func fetchChain(w http.ResponseWriter, r *http.Request) {
	jsontxt, err := json.Marshal(servBlockchain.Chain)
	if err != nil {
		log.Println("Error: " + err.Error())
		return
	}
	w.Write(jsontxt)
}

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
	decoder := json.NewDecoder(r.Body)
	var newPoll Poll
	err := decoder.Decode(&newPoll)
	if err != nil {
		log.Println("New Poll Data Parse Error:", err)
		w.Write([]byte("#"))
	}
	fmt.Println("Adding new poll:", newPoll)
	pollPass := shaHash(strconv.FormatInt(time.Now().UTC().UnixNano(), 10))
	newPoll.PassHash = shaHash(pollPass)
	pollID := createPoll(newPoll)
	w.Write([]byte("/admin/" + strconv.Itoa(pollID) + "/" + pollPass))
}
