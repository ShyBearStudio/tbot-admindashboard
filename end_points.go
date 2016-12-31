package main

import (
	"net/http"

	"github.com/ShyBearStudio/tbot-admindashboard/data"
	"github.com/ShyBearStudio/tbot-admindashboard/logger"
)

type routeFunc func(http.ResponseWriter, *http.Request)

type endPointDesc struct {
	pattern       string
	loginRequired bool
	roleRouting   map[data.UserRoleType]routeFunc
}

func (endPoint *endPointDesc) addRoute(
	role data.UserRoleType, route routeFunc) *endPointDesc {
	endPoint.roleRouting[role] = route
	return endPoint
}

func (endPoint *endPointDesc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !endPoint.loginRequired {
		endPoint.roleRouting[data.AllRoles](w, r)
		return
	}
	sess, err := session(r)
	if err != nil {
		redirect(w, r, endPoints.login)
		return
	}

	user, err := sess.User()
	if err != nil {
		logger.Errorln("Cannot file by session", err)
	}
	endPoint.roleRouting[user.Role](w, r)
}

type endPointsList struct {
	index endPointDesc
	login endPointDesc
}

var endPoints endPointsList

func init() {
	endPoints = endPointsList{}
	// Index end point
	endPoints.index = endPointDesc{pattern: "/", loginRequired: true, roleRouting: make(map[data.UserRoleType]routeFunc)}
	endPoints.index.addRoute(data.AllRoles, index)
	// Login end point
	endPoints.login = endPointDesc{pattern: "/login", loginRequired: false, roleRouting: make(map[data.UserRoleType]routeFunc)}
	endPoints.login.addRoute(data.AllRoles, login)
}

func redirect(w http.ResponseWriter, r *http.Request, endPoint endPointDesc) {
	http.Redirect(w, r, endPoint.pattern, 302)
}
