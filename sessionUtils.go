package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ShyBearStudio/tbot-admindashboard/data"
)

// Checks if the user is logged in and has a session, if not err is not nil
func session(r *http.Request) (sess data.Session, err error) {
	_ = "breakpoint"
	cookie, err := sessionCookie(r)
	fmt.Printf("cookie: %v\n", cookie)
	if err == nil {
		sess = data.Session{Uuid: cookie.Value}
		if ok, _ := sess.Check(); !ok {
			err = errors.New("Invalid session")
		}
	}
	return
}

func sessionCookie(r *http.Request) (cookie *http.Cookie, err error) {
	cookie, err = r.Cookie("_cookie")
	return
}
