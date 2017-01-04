package main

import (
	"fmt"
	"net/http"

	"github.com/ShyBearStudio/tbot-admindashboard/data"
	"github.com/ShyBearStudio/tbot-admindashboard/logger"
)

type endPoint struct {
	pattern string
	// method -> handler
	routing map[string]http.Handler
}

func newEndPoint(pattern string) endPoint {
	ep := endPoint{}
	ep.pattern = pattern
	ep.routing = make(map[string]http.Handler)
	return ep
}

type endPointId int

func (id endPointId) String() string {
	return fmt.Sprintf("%d", int(id))
}

const (
	notFountEndPoint = iota + 1
	indexEndPoint
	loginEndPoint
	authEndPoint
	createAccountEndPoint
	usersEndPoint
	projectsEndPoint
	echoProjectEndPoint
)

var endPoints = make(map[endPointId]endPoint)

type handlerFunc func(http.ResponseWriter, *http.Request)

type authEndPointHandler struct {
	handler handlerFunc
}

func newAuthEndPointHandler() *authEndPointHandler {
	h := new(authEndPointHandler)
	return h
}

func (h *authEndPointHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = "breakpoint"
	h.handler(w, r)
}

type roleBasedEndPointHandler struct {
	handler        map[data.UserRoleType]handlerFunc
	defaultHandler handlerFunc
}

func newRoleBasedEndPointHandler() *roleBasedEndPointHandler {
	h := new(roleBasedEndPointHandler)
	h.handler = make(map[data.UserRoleType]handlerFunc)
	return h
}

func (h *roleBasedEndPointHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = "breakpoint"
	var url = r.URL.Path[1:]
	sess, err := session(r)
	if err != nil {
		logger.Warningf("Unauthorized access to '%s'", url)
		redirect(loginEndPoint, w, r)
		return
	}

	user, err := sess.User()
	if err != nil {
		logger.Errorln("Cannot find user by session", err)
	}

	handler, ok := h.handler[user.Role]
	if !ok {
		if h.defaultHandler != nil {
			h.defaultHandler(w, r)
			return
		}

		redirect(notFountEndPoint, w, r)
		return
	}
	handler(w, r)
}

func redirect(id endPointId, w http.ResponseWriter, r *http.Request) {
	ep, ok := endPoints[id]
	if !ok {
		logger.Errorf("There is no endpoint with id: '%d'")
		return
	}
	http.Redirect(w, r, ep.pattern, 302)
}

func init() {
	endPoints[notFountEndPoint] = initNotFoundEndPoint()
	endPoints[indexEndPoint] = initIndexEndPoint()
	endPoints[loginEndPoint] = initLoginEndPoint()
	endPoints[authEndPoint] = initAuthEndPoint()
	endPoints[createAccountEndPoint] = initCreateAccountEndPoint()
	endPoints[usersEndPoint] = initUsersEndPoint()
	endPoints[projectsEndPoint] = initProjectsEndPoint()
	endPoints[echoProjectEndPoint] = initEchoProjectEndPoint()
}

func initNotFoundEndPoint() endPoint {
	ep := newEndPoint("/not_found")
	getRouting := newAuthEndPointHandler()
	getRouting.handler = notFount
	return ep
}

func initIndexEndPoint() endPoint {
	ep := newEndPoint("/")
	getRouting := newRoleBasedEndPointHandler()
	getRouting.defaultHandler = index
	ep.routing[http.MethodGet] = getRouting
	return ep
}

func initLoginEndPoint() endPoint {
	ep := newEndPoint("/login")
	getRouting := newAuthEndPointHandler()
	getRouting.handler = login
	ep.routing[http.MethodGet] = getRouting
	return ep
}

func initAuthEndPoint() endPoint {
	ep := newEndPoint("/authenticate")
	postRouting := newAuthEndPointHandler()
	postRouting.handler = authenticate
	ep.routing[http.MethodPost] = postRouting
	return ep
}

func initCreateAccountEndPoint() endPoint {
	ep := newEndPoint("/create_account")
	getRouting := newRoleBasedEndPointHandler()
	getRouting.handler[data.AdminRole] = createAccountGet
	ep.routing[http.MethodGet] = getRouting
	postRouting := newRoleBasedEndPointHandler()
	postRouting.handler[data.AdminRole] = createAccountPost
	ep.routing[http.MethodPost] = postRouting
	return ep
}

func initUsersEndPoint() endPoint {
	ep := newEndPoint("/users")
	getRouting := newRoleBasedEndPointHandler()
	getRouting.handler[data.AdminRole] = users
	ep.routing[http.MethodGet] = getRouting
	return ep
}

func initProjectsEndPoint() endPoint {
	ep := newEndPoint("/projects")
	getRouting := newRoleBasedEndPointHandler()
	getRouting.handler[data.AdminRole], getRouting.handler[data.UserRole] = projects, projects
	ep.routing[http.MethodGet] = getRouting
	return ep
}

func initEchoProjectEndPoint() endPoint {
	ep := newEndPoint("/projects/echo")
	getRouting := newRoleBasedEndPointHandler()
	getRouting.handler[data.AdminRole], getRouting.handler[data.UserRole] = echoProject, echoProject
	ep.routing[http.MethodGet] = getRouting
	return ep
}
