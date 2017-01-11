package main

import (
	"fmt"
	"net/http"

	"github.com/ShyBearStudio/tbot-admindashboard/logger"
	"github.com/ShyBearStudio/tbot-admindashboard/models"
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
	env            *environment
	handler        map[models.UserRoleType]handlerFunc
	defaultHandler handlerFunc
}

func newRoleBasedEndPointHandler(env *environment) *roleBasedEndPointHandler {
	h := new(roleBasedEndPointHandler)
	h.env = env
	h.handler = make(map[models.UserRoleType]handlerFunc)
	return h
}

func (h *roleBasedEndPointHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = "breakpoint"
	var url = r.URL.Path[1:]
	sess, err := session(h.env, r)
	if err != nil {
		logger.Warningf("Unauthorized access to '%s'", url)
		h.env.redirect(loginEndPoint, w, r)
		return
	}

	user, err := h.env.db.User(&sess)
	if err != nil {
		logger.Errorln("Cannot find user by session", err)
	}

	handler, ok := h.handler[user.Role]
	if !ok {
		if h.defaultHandler != nil {
			h.defaultHandler(w, r)
			return
		}

		h.env.redirect(notFountEndPoint, w, r)
		return
	}
	handler(w, r)
}

func (env *environment) redirect(id endPointId, w http.ResponseWriter, r *http.Request) {
	ep, ok := env.endPoints[id]
	if !ok {
		logger.Errorf("There is no endpoint with id: '%d'")
		return
	}
	http.Redirect(w, r, ep.pattern, 302)
}

//func init() {
func newEndpoins(env *environment) map[endPointId]endPoint {
	endPoints := make(map[endPointId]endPoint)

	endPoints[notFountEndPoint] = initNotFoundEndPoint(env)
	endPoints[indexEndPoint] = initIndexEndPoint(env)
	endPoints[loginEndPoint] = initLoginEndPoint(env)
	endPoints[authEndPoint] = initAuthEndPoint(env)
	endPoints[createAccountEndPoint] = initCreateAccountEndPoint(env)
	endPoints[usersEndPoint] = initUsersEndPoint(env)
	endPoints[projectsEndPoint] = initProjectsEndPoint(env)
	endPoints[echoProjectEndPoint] = initEchoProjectEndPoint(env)

	return endPoints
}

func initNotFoundEndPoint(env *environment) endPoint {
	ep := newEndPoint("/not_found")
	getRouting := newAuthEndPointHandler()
	getRouting.handler = env.notFount
	return ep
}

func initIndexEndPoint(env *environment) endPoint {
	ep := newEndPoint("/")
	getRouting := newRoleBasedEndPointHandler(env)
	getRouting.defaultHandler = env.index
	ep.routing[http.MethodGet] = getRouting
	return ep
}

func initLoginEndPoint(env *environment) endPoint {
	ep := newEndPoint("/login")
	getRouting := newAuthEndPointHandler()
	getRouting.handler = env.login
	ep.routing[http.MethodGet] = getRouting
	return ep
}

func initAuthEndPoint(env *environment) endPoint {
	ep := newEndPoint("/authenticate")
	postRouting := newAuthEndPointHandler()
	postRouting.handler = env.authenticate
	ep.routing[http.MethodPost] = postRouting
	return ep
}

func initCreateAccountEndPoint(env *environment) endPoint {
	ep := newEndPoint("/create_account")
	getRouting := newRoleBasedEndPointHandler(env)
	getRouting.handler[models.AdminRole] = env.createAccountGet
	ep.routing[http.MethodGet] = getRouting
	postRouting := newRoleBasedEndPointHandler(env)
	postRouting.handler[models.AdminRole] = env.createAccountPost
	ep.routing[http.MethodPost] = postRouting
	return ep
}

func initUsersEndPoint(env *environment) endPoint {
	ep := newEndPoint("/users")
	getRouting := newRoleBasedEndPointHandler(env)
	getRouting.handler[models.AdminRole] = env.users
	ep.routing[http.MethodGet] = getRouting
	return ep
}

func initProjectsEndPoint(env *environment) endPoint {
	ep := newEndPoint("/projects")
	getRouting := newRoleBasedEndPointHandler(env)
	getRouting.handler[models.AdminRole] = env.projects
	getRouting.handler[models.UserRole] = env.projects
	ep.routing[http.MethodGet] = getRouting
	return ep
}

func initEchoProjectEndPoint(env *environment) endPoint {
	ep := newEndPoint("/projects/echo")
	getRouting := newRoleBasedEndPointHandler(env)
	getRouting.handler[models.AdminRole] = env.echoProject
	getRouting.handler[models.UserRole] = env.echoProject
	ep.routing[http.MethodGet] = getRouting
	return ep
}
