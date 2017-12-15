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
	"github.com/subosito/gotenv"
)

var servBlockchain = BlockChain{}.start()
var nodeIdentifier = uuid.NewV4().String()

func init() {
	gotenv.Load()
	mongoLoad()
	getPolls()
}
func main() {
	// Process handlebars templates
	template.AddEngine(handlebars.New()).Directory("./views", ".hbs")
	err := template.Load()
	if err != nil {
		panic("While parsing the template files: " + err.Error())
	}
	//Request routing
	router := mux.NewRouter()
	//Resources
	router.PathPrefix("/res").Handler(http.StripPrefix("/res", http.FileServer(http.Dir("res/"))))
	//UI routing
	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/vote/{pollID:[0-9]+}", pollPage).Methods("GET")
	router.HandleFunc("/newpoll", newpollPage).Methods("GET")
	//Main blockchain routing
	router.HandleFunc("/castvote", voteAPI).Methods("GET")
	//Start the engines
	portPtr := flag.String("p", "3333", "Server Port")
	flag.Parse()
	port := ":" + *portPtr
	color.Green("Starting server on port: %s", port[1:])
	log.Fatal(http.ListenAndServe(port, router))
}
