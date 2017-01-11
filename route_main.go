package main

import (
	"fmt"
	"net/http"

	"github.com/ShyBearStudio/tbot-admindashboard/logger"
)

func (env *environment) index(w http.ResponseWriter, r *http.Request) {
	_ = "breakpoint"
	fmt.Fprint(w, "access granted!\n")
	fmt.Fprint(w, "Hello World!")
}

func (env *environment) notFount(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Not found Page (404)!")
}

func (env *environment) users(w http.ResponseWriter, r *http.Request) {
	_ = "breakpoint"
	users, err := env.db.Users()
	if err != nil {
		logger.Errorln("Cannot get all users", err)
	}
	generateHTML(w, &users, "page.layout", "users")
}
