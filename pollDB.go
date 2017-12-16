package main

import (
	"fmt"
	"log"
	"os"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Poll struct {
	ID       int
	Name     string
	Options  []string
	PassHash string
}

var mongoURI string
var mongoDB *mgo.Session

func mongoLoad() {
	mongoURI = os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("No MongoDB URI supplied in .env config file!")
	}
	var err error
	mongoDB, err = mgo.Dial(mongoURI)
	if err != nil {
		log.Fatal("Failed to connect to provided MongoDB URI:\n", err)
	}
}
func createPoll(pollDoc Poll) int {
	mongoSesh := mongoDB.Copy()
	defer mongoSesh.Close()
	pollCt, _ := mongoSesh.DB("blockvote").C("polls").Find(bson.M{}).Count()
	pollDoc.ID = pollCt
	err := mongoSesh.DB("blockvote").C("polls").Insert(pollDoc)
	if err != nil {
		fmt.Println("Failure to insert poll document:\n", err)
	}
	pollDictMutex.Lock()
	pollDict[pollCt] = pollDoc
	pollDictMutex.Unlock()
	return pollCt
}
func getPolls() []Poll {
	mongoSesh := mongoDB.Copy()
	defer mongoSesh.Close()
	var pollList []Poll
	mongoSesh.DB("blockvote").C("polls").Find(bson.M{}).All(&pollList)
	fmt.Println(pollList)
	return pollList
}
func pollListToDict(pollList []Poll) map[int]Poll {
	pollDict := make(map[int]Poll, len(pollList))
	for _, pollElem := range pollList {
		pollDict[pollElem.ID] = pollElem
	}
	return pollDict
}
