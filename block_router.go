package main

import (
	"encoding/json"
	"log"
	"net/http"
)

//Main blockchain routing
func mine(w http.ResponseWriter, r *http.Request) {
	lastBlock := servBlockchain.lastBlock()
	lastProof := lastBlock.Proof
	proof := servBlockchain.proofOfWork(lastProof)
	//servBlockchain.newTransaction("0", nodeIdentifier, 1)
	previousHash := lastBlock.hash()
	newBlock := servBlockchain.newBlock(proof, previousHash)
	response := map[string]interface{}{
		"message":       "New Block Forged",
		"index":         newBlock.Index,
		"transactions":  newBlock.Transactions,
		"proof":         newBlock.Proof,
		"previous_hash": newBlock.PreviousHash,
	}
	jsontxt, err := json.Marshal(response)
	if err != nil {
		log.Println("Error: " + err.Error())
		return
	}
	w.Write(jsontxt)
}

func fetchChain(w http.ResponseWriter, r *http.Request) {
	jsontxt, err := json.Marshal(servBlockchain.Chain)
	if err != nil {
		log.Println("Error: " + err.Error())
		return
	}
	w.Write(jsontxt)
}

func newTransaction(w http.ResponseWriter, r *http.Request) {
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
	servBlockchain.newTransaction(Transaction{voterid[0], vote[0]})
	w.Write([]byte("Vote sucessfully recorded."))
}
