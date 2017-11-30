package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

//Main blockchain routing
func mine(w http.ResponseWriter, r *http.Request) {
	lastBlock := servBlockchain.lastBlock()
	lastProof := lastBlock.Proof
	proof := servBlockchain.proofOfWork(lastProof)
	servBlockchain.newTransaction("0", nodeIdentifier, 1)
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
	sender, ok := r.URL.Query()["sender"]
	if !ok || len(sender) < 1 {
		http.Error(w, "Err 400: URL Param 'sender' is missing", 400)
		return
	}
	receiver, ok := r.URL.Query()["receiver"]
	if !ok || len(receiver) < 1 {
		http.Error(w, "Err 400: URL Param 'receiver' is missing", 400)
		return
	}
	amount, ok := r.URL.Query()["amount"]
	if !ok || len(amount) < 1 {
		http.Error(w, "Err 400: URL Param 'amount' is missing", 400)
		return
	}
	amountNum, err := strconv.Atoi(amount[0])
	if err != nil {
		http.Error(w, "Err 400: URL Param 'amount' is invalid", 400)
		return
	}
	index := servBlockchain.newTransaction(sender[0], receiver[0], amountNum)
	w.Write([]byte("Transaction will be added to Block: " + strconv.Itoa(index)))
}
