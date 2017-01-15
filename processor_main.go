package main

import (
	"fmt"
	"net/http"

	"github.com/ShyBearStudio/tbot-admindashboard/logger"
)

func (env *environment) index(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	_ = "breakpoint"
	fmt.Fprint(w, "access granted!\n")
	fmt.Fprint(w, "Hello World!")
	generateHTML(w, nil, "navbar", "page.layout")
	return nil, nil
}

func (env *environment) notFoundBody(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Not found Page (404)!")
}

func (env *environment) notFound(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	env.notFoundBody(w, r)
	return nil, nil
}

func (env *environment) users(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	_ = "breakpoint"
	users, err := env.db.Users()
	if err != nil {
		logger.Errorln("Cannot get all users", err)
	}
	return users, err
}

func (env *environment) emptyProcessor(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, nil
}
