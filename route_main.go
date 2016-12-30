package main

import (
	"fmt"
	"net/http"

	_ "github.com/ShyBearStudio/tbot-admindashboard/data"
)

/*
type routingFunc func(http.ResponseWriter, *http.Request)

var endpointPrivMap map[string]map[data.UserRoleType]routingFunc

func init() {
	endpointPrivMap["index"][data.AdminRole] =
	endpointPrivMap["index"][data.UserRole] = index
}*/

func index(w http.ResponseWriter, r *http.Request) {
	_ = "breakpoint"
	_, err := session(r)
	if err != nil {
		login(w, r)
		return
		//fmt.Fprint(w, "need to log in!\n")
	} else {
		fmt.Fprint(w, "access granted!\n")
	}
	fmt.Fprint(w, "Hello, World!")
}

/*
func f(endPoint string, w http.ResponseWriter, r *http.Request) {
	sess, err := session(r)
	if err != nil {
		login(w, r)
	}
	user, err := sess.User()
	role := user.Role
	// endpoint + role -> route function
}
*/
