package main

import (
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	template "github.com/kataras/go-template"
	"github.com/kataras/go-template/handlebars"
	uuid "github.com/satori/go.uuid"
	"github.com/subosito/gotenv"
)

var servBlockchain = BlockChain{}.start()
var nodeIdentifier = uuid.NewV4().String()
var pollDict map[int]Poll
var pollDictMutex sync.Mutex

func init() {
	gotenv.Load()
	mongoLoad()
	pollDictMutex.Lock()
	pollDict = pollListToDict(getPolls())
	pollDictMutex.Unlock()
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
	router.HandleFunc("/vote/{pollID:[0-9]+}/:pollPass", pollAdminPage).Methods("GET")
	router.HandleFunc("/vote/{pollID:[0-9]+}", pollPage).Methods("GET")
	router.HandleFunc("/newpoll", newpollPage).Methods("GET")
	//API routing
	router.HandleFunc("/api/registerpoll", newPollAPI).Methods("POST")
	router.HandleFunc("/api/chain", fetchChain).Methods("GET")
	router.HandleFunc("/api/castvote", voteAPI).Methods("GET")
	//Start the engines
	portPtr := flag.String("p", "3333", "Server Port")
	flag.Parse()
	port := ":" + *portPtr
	color.Green("Starting server on port: %s", port[1:])
	log.Fatal(http.ListenAndServe(port, router))
}
