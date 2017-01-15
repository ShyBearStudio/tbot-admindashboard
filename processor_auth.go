package main

import (
	"net/http"

	"github.com/ShyBearStudio/tbot-admindashboard/dbutils"
	"github.com/ShyBearStudio/tbot-admindashboard/logger"
	"github.com/ShyBearStudio/tbot-admindashboard/models"
)

func (env *environment) login(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, nil
}

// POST /authenticate
// Authenticate the user given the email and password
func (env *environment) authenticate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	_ = "breakpoint"
	err := r.ParseForm()
	if err != nil {
		logger.Errorln("Cannot parse form", err)
	}
	var email = r.PostFormValue("email")
	user, err := env.db.UserByEmail(email)
	if err != nil {
		logger.Errorf("Cannot find user with email '%s': %v", email, err)
		env.redirect(loginEndPoint, w, r)
	}
	logger.Tracef("Try to authenticate user '%v'.", user)
	if user.Password == dbutils.Encrypt(r.PostFormValue("password")) {
		session, err := env.db.CreateSession(&user)
		if err != nil {
			logger.Errorln("Cannot create session", err)
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		env.redirect(indexEndPoint, w, r)
	} else {
		logger.Errorf("Incorrect password for user '%v'", user)
		env.redirect(loginEndPoint, w, r)
	}
	return nil, err
}

func (env *environment) createAccountPost(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := r.ParseForm()
	if err != nil {
		logger.Errorln("Cannot parse create acccount form", err)
	}
	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	// Right away we allow only user creation since we've not got to the bottom of
	// vendor role and do not have final design
	role := models.UserRole
	logger.Tracef("Creating user. Name: '%s' | Email: '%s' | Password: '%s' | Role: '%s'", name, email, password, role)
	userId, err := env.db.AddUser(name, email, password, role)
	if err != nil {
		logger.Errorln("Cannot create account", err)
	}
	logger.Tracef("Account with '%d' created", userId)
	env.redirect(indexEndPoint, w, r)
	return nil, nil
}
