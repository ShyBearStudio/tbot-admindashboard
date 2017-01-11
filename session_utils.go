package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ShyBearStudio/tbot-admindashboard/models"
)

// Checks if the user is logged in and has a session, if not err is not nil
func session(env *environment, r *http.Request) (sess models.Session, err error) {
	_ = "breakpoint"
	cookie, err := sessionCookie(r)
	fmt.Printf("cookie: %v\n", cookie)
	if err == nil {
		sess = models.Session{Uuid: cookie.Value}
		if ok, _ := env.db.CheckSession(&sess); !ok {
			err = errors.New("Invalid session")
		}
	}
	return
}

func sessionCookie(r *http.Request) (cookie *http.Cookie, err error) {
	cookie, err = r.Cookie("_cookie")
	return
}
