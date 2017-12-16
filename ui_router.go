package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	template "github.com/kataras/go-template"
)

var renderOpts = map[string]interface{}{"layout": "layouts/base.hbs"}

func renderPage(templateName string, pageData map[string]interface{}, w http.ResponseWriter) {
	err := template.ExecuteWriter(w, templateName, pageData, renderOpts) // yes you can pass simple maps instead of structs
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
func homePage(w http.ResponseWriter, r *http.Request) {
	err := template.ExecuteWriter(w, "home.hbs", map[string]interface{}{}, renderOpts) // yes you can pass simple maps instead of structs
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
func pollAdminPage(w http.ResponseWriter, r *http.Request) {
	urlparams := mux.Vars(r)
	pollID := urlparams["pollID"]
	pollStringID, _ := strconv.Atoi(pollID)
	pollPass := urlparams["pollPass"]
	pollObj, ok := pollDict[pollStringID]
	if ok {
		if shaHash(pollPass) == pollObj.PassHash {
			renderPage("vote.hbs", map[string]interface{}{"pollID": pollID, "voteopts": pollObj.Options}, w)
		} else {
			redirectPage(w, "/", "Invalid Credentials!", "1.5")
		}
	} else {
		redirectPage(w, "/", "Cannot find poll with this ID!", "1.5")
	}
}
func pollPage(w http.ResponseWriter, r *http.Request) {
	urlparams := mux.Vars(r)
	pollID := urlparams["pollID"]
	pollStringID, _ := strconv.Atoi(pollID)
	pollObj, ok := pollDict[pollStringID]
	if ok {
		renderPage("vote.hbs", map[string]interface{}{"pollID": pollID, "voteopts": pollObj.Options}, w)
		err := template.ExecuteWriter(w, "vote.hbs", map[string]interface{}{"pollID": pollID, "voteopts": pollObj.Options}, renderOpts) // yes you can pass simple maps instead of structs
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	} else {
		redirectPage(w, "/", "Cannot find poll with this ID!", "1.5")
	}
}
func newpollPage(w http.ResponseWriter, r *http.Request) {
	err := template.ExecuteWriter(w, "newpoll.hbs", map[string]interface{}{}, renderOpts) // yes you can pass simple maps instead of structs
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

//Helper function
func redirectPage(w http.ResponseWriter, redirUrl string, redirMsg string, redirTime string) {
	err := template.ExecuteWriter(w, "redirect.hbs", map[string]interface{}{
		"redir_url":  redirUrl,
		"redir_msg":  redirMsg,
		"redir_time": redirTime,
	}, renderOpts) // yes you can pass simple maps instead of structs
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
