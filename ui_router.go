package main

import (
	"net/http"

	"github.com/gorilla/mux"
	template "github.com/kataras/go-template"
)

var renderOpts = map[string]interface{}{"layout": "layouts/base.hbs"}

func homePage(w http.ResponseWriter, r *http.Request) {
	err := template.ExecuteWriter(w, "home.hbs", map[string]interface{}{}, renderOpts) // yes you can pass simple maps instead of structs
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
func pollPage(w http.ResponseWriter, r *http.Request) {
	urlparams := mux.Vars(r)
	pollID := urlparams["pollID"]
	//if pollID != "" &&
	voteOptions := []string{"Jordan Poles", "Maria D'lorio"}
	err := template.ExecuteWriter(w, "vote.hbs", map[string]interface{}{"pollID": pollID, "voteopts": voteOptions}, renderOpts) // yes you can pass simple maps instead of structs
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
