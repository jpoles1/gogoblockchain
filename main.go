package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

var servBlockchain BlockChain = BlockChain{}.start()
var nodeIdentifier = uuid.NewV4().String()

func main() {
	var jsonChain []Block
	x := `[{"Index":1,"Timestamp":1511930004931818700,"Transactions":null,"Proof":100,"PreviousHash":"1"},{"Index":2,"Timestamp":1511930123436683600,"Transactions":[{"Sender":"0","Receipient":"49d8edbb-cb02-4441-8ff0-8a0a8550aa89","Amount":1}],"Proof":3707,"PreviousHash":"da6g0TFhyEIlJuNfIylwwfHCpO2vRy2aVU2GMHk7KuI="}]`
	if err := json.Unmarshal([]byte(x), &jsonChain); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(jsonChain[1])

	}
	//
	router := mux.NewRouter()
	//Main blockchain routing
	router.HandleFunc("/transactions/new", newTransaction).Methods("GET")
	router.HandleFunc("/chain", fetchChain).Methods("GET")
	router.HandleFunc("/mine", mine).Methods("GET")
	//Multinode routing
	router.HandleFunc("/nodes/register", registerNode).Methods("GET")
	router.HandleFunc("/nodes/resolve", resolveNode).Methods("GET")
	//Start the engines
	port := ":3333"
	color.Green("Starting server on port: %s", port[1:])
	log.Fatal(http.ListenAndServe(port, router))
}
