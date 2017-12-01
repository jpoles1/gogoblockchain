package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	template "github.com/kataras/go-template"
	"github.com/kataras/go-template/handlebars"
	uuid "github.com/satori/go.uuid"
)

var servBlockchain = BlockChain{}.start()
var nodeIdentifier = uuid.NewV4().String()

func main() {
	// Process handlebars templates
	template.AddEngine(handlebars.New()).Directory("./views", ".hbs")
	err := template.Load()
	if err != nil {
		panic("While parsing the template files: " + err.Error())
	}
	//Request routing

	router := mux.NewRouter()
	//UI routing
	router.HandleFunc("/", votePage).Methods("GET")
	//Main blockchain routing
	router.HandleFunc("/transactions/new", newTransaction).Methods("GET")
	router.HandleFunc("/chain", fetchChain).Methods("GET")
	router.HandleFunc("/mine", mine).Methods("GET")
	//Multinode routing
	router.HandleFunc("/nodes/register", registerNode).Methods("GET")
	router.HandleFunc("/nodes/resolve", resolveNode).Methods("GET")
	//Start the engines
	portPtr := flag.String("p", "3333", "Server Port")
	flag.Parse()
	port := ":" + *portPtr
	color.Green("Starting server on port: %s", port[1:])
	log.Fatal(http.ListenAndServe(port, router))
}
