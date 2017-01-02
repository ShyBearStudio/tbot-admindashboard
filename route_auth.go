package main

import (
	"net/http"

	"github.com/ShyBearStudio/tbot-admindashboard/data"
	"github.com/ShyBearStudio/tbot-admindashboard/logger"
)

func login(w http.ResponseWriter, r *http.Request) {
	t := parseTemplateFiles("login.layout", "login")
	t.Execute(w, nil)
}

// POST /authenticate
// Authenticate the user given the email and password
func authenticate(w http.ResponseWriter, r *http.Request) {
	_ = "breakpoint"
	err := r.ParseForm()
	if err != nil {
		logger.Errorln("Cannot parse form", err)
	}
	var email = r.PostFormValue("email")
	user, err := data.UserByEmail(email)
	if err != nil {
		logger.Errorf("Cannot find user with email '%s': %v", email, err)
		redirect(loginEndPoint, w, r)
	}
	logger.Tracef("Try to authenticate user '%v'.", user)
	if user.Password == data.Encrypt(r.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			logger.Errorln("Cannot create session", err)
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		redirect(indexEndPoint, w, r)
	} else {
		logger.Errorf("Incorrect password for user '%v'", user)
		redirect(loginEndPoint, w, r)
	}
}

func createAccountGet(w http.ResponseWriter, r *http.Request) {
	_ = "breakpoint"
	generateHTML(w, nil, "login.layout", "create_account")
}

func createAccountPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		logger.Errorln("Cannot parse create acccount form", err)
	}
	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	// Right away we allow only user creation since we've not got to the bottom of
	// vendor role and do not have final design
	role := data.UserRole
	logger.Tracef("Creating user. Name: '%s' | Email: '%s' | Password: '%s' | Role: '%s'", name, email, password, role)
	userId, err := data.CreateUser(name, email, password, role)
	if err != nil {
		logger.Errorln("Cannot create account", err)
	}
	logger.Tracef("Account with '%d' created", userId)
	redirect(indexEndPoint, w, r)
}
