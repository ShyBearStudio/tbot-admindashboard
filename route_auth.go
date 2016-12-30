package main

import (
	"net/http"

	_ "github.com/ShyBearStudio/tbot-admindashboard/data"
)

func login(w http.ResponseWriter, r *http.Request) {
	t := parseTemplateFiles("login.layout", "login")
	t.Execute(w, nil)
}
