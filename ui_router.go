package main

import (
	"net/http"

	template "github.com/kataras/go-template"
)

var renderOpts = map[string]interface{}{"layout": "layouts/base.hbs"}

func votePage(w http.ResponseWriter, r *http.Request) {
	voteOptions := []string{"Jordan Poles", "Maria D'lorio"}
	err := template.ExecuteWriter(w, "vote.hbs", map[string]interface{}{"voteopts": voteOptions}, renderOpts) // yes you can pass simple maps instead of structs
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
