package main

import (
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

var servBlockchain BlockChain = BlockChain{}.start()
var nodeIdentifier = uuid.NewV4().String()

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/transactions/new", newTransaction).Methods("GET")
	router.HandleFunc("/chain", fetchChain).Methods("GET")
	router.HandleFunc("/mine", mine).Methods("GET")
	port := ":3333"
	color.Green("Starting server on port: %s", port[1:])
	log.Fatal(http.ListenAndServe(port, router))
}
