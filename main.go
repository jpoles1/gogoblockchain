package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	template "github.com/kataras/go-template"
	"github.com/kataras/go-template/handlebars"
	uuid "github.com/satori/go.uuid"
	"github.com/subosito/gotenv"
)

var servBlockchain = BlockChain{}
var nodeIdentifier = uuid.NewV4().String()
var pollDict map[int]Poll
var pollDictMutex sync.Mutex

func init() {
	gotenv.Load()
	initEmail()
	mongoLoad()
	//Load list of polls
	pollDictMutex.Lock()
	pollDict = pollListToDict(getPolls())
	pollDictMutex.Unlock()
	//Load backup of blockchain
	blockchainData, err := ioutil.ReadFile("data/blockchain.json")
	if err != nil {
		fmt.Println("Cannot load blockchain backup:", err)
		servBlockchain = servBlockchain.start()
		servBlockchain.newBlock(100, "1")
	} else {
		json.Unmarshal(blockchainData, &servBlockchain)
		servBlockchain = servBlockchain.start()
		if len(servBlockchain.Chain) < 1 {
			servBlockchain.newBlock(100, "1")
		}
	}
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
	router.HandleFunc("/admin/{pollID:[0-9]+}/{pollPass}", pollAdminPage).Methods("GET")
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
	//Handling system signals
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println("\nReceived Command:", sig)
		done <- true
	}()
	go http.ListenAndServe(port, router)
	<-done
	fmt.Println("Terminating BlockVote Server...")
	jsontxt, err := json.Marshal(servBlockchain)
	if err != nil {
		log.Println("BlockChain Save Error: " + err.Error())
		return
	}
	err = ioutil.WriteFile("data/blockchain.json", jsontxt, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Saved blockchain. Exiting gracefully...\n")
}
