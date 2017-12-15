package main

import (
	"net/http"
	"strconv"

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
	pollStringID, _ := strconv.Atoi(pollID)
	pollObj, ok := pollDict[pollStringID]
	if ok {
		err := template.ExecuteWriter(w, "vote.hbs", map[string]interface{}{"pollID": pollID, "voteopts": pollObj.Options}, renderOpts) // yes you can pass simple maps instead of structs
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	} else {
		err := template.ExecuteWriter(w, "redirect.hbs", map[string]interface{}{
			"redir_msg": "Cannot find poll with this ID!",
			"redir_url": "/",
		}, renderOpts) // yes you can pass simple maps instead of structs
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	}
}
func newpollPage(w http.ResponseWriter, r *http.Request) {
	err := template.ExecuteWriter(w, "newpoll.hbs", map[string]interface{}{}, renderOpts) // yes you can pass simple maps instead of structs
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
