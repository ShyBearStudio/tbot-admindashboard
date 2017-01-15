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
	logoutEndPoint
	authEndPoint
	createAccountEndPoint
	usersEndPoint
	projectsEndPoint
	echoProjectEndPoint
)

type processorFunc func(http.ResponseWriter, *http.Request) (interface{}, error)

type baseEndPointHandler struct {
	env           *environment
	htmlFileNames []string
}

type authEndPointHandler struct {
	baseEndPointHandler
	processor processorFunc
}

func newAuthEndPointHandler(
	env *environment, processor processorFunc, htmls ...string) *authEndPointHandler {
	h := new(authEndPointHandler)
	h.env = env
	h.processor = processor
	h.htmlFileNames = htmls
	return h
}

func (h *authEndPointHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = "breakpoint"
	data, err := h.processor(w, r)
	if err != nil {
		logger.Errorln(err)
	}
	generateHTML(w, data, h.htmlFileNames...)
	//h.handler(w, r)
}

type roleBasedEndPointHandler struct {
	baseEndPointHandler
	processors     map[models.UserRoleType]processorFunc
	defaultHandler processorFunc
}

func newRoleBasedEndPointHandler(env *environment, htmls ...string) *roleBasedEndPointHandler {
	h := new(roleBasedEndPointHandler)
	h.env = env
	h.processors = make(map[models.UserRoleType]processorFunc)
	h.htmlFileNames = htmls
	return h
}

func (h *roleBasedEndPointHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = "breakpoint"
	user, err := h.User(r)
	if err != nil {
		logger.Errorln("Cannot find user by session", err)
		h.env.redirect(loginEndPoint, w, r)
		return
	}

	processor := h.RouteProcessor(user.Role)
	if processor == nil {
		h.env.redirect(notFountEndPoint, w, r)
		return
	}

	data, err := processor(w, r)
	if err != nil {
		logger.Errorln(err)
	}
	generateHTML(w, data, h.htmlFileNames...)
}

func (h *roleBasedEndPointHandler) User(r *http.Request) (user *models.User, err error) {
	sess, err := session(h.env, r)
	if err != nil {
		logger.Warningf("Unauthorized access to '%s'", r.URL.Path[1:])
		return nil, err
	}
	user, err = h.env.db.User(sess)
	return
}

func (h *roleBasedEndPointHandler) RouteProcessor(role models.UserRoleType) processorFunc {
	processor, ok := h.processors[role]
	if !ok && h.defaultHandler != nil {
		processor = h.defaultHandler
	}
	return processor
}

func (env *environment) redirect(id endPointId, w http.ResponseWriter, r *http.Request) {
	ep, ok := env.endPoints[id]
	if !ok {
		logger.Errorf("There is no endpoint with id: '%d'")
		return
	}
	http.Redirect(w, r, ep.pattern, 302)
}

func newEndpoins(env *environment) map[endPointId]endPoint {
	endPoints := make(map[endPointId]endPoint)

	endPoints[notFountEndPoint] = initNotFoundEndPoint(env)
	endPoints[indexEndPoint] = initIndexEndPoint(env)
	endPoints[loginEndPoint] = initLoginEndPoint(env)
	endPoints[logoutEndPoint] = initLogoutEndPoint(env)
	endPoints[authEndPoint] = initAuthEndPoint(env)
	endPoints[createAccountEndPoint] = initCreateAccountEndPoint(env)
	endPoints[usersEndPoint] = initUsersEndPoint(env)
	endPoints[projectsEndPoint] = initProjectsEndPoint(env)
	endPoints[echoProjectEndPoint] = initEchoProjectEndPoint(env)

	return endPoints
}

func initNotFoundEndPoint(env *environment) endPoint {
	ep := newEndPoint("/not_found")
	ep.routing[http.MethodGet] = newAuthEndPointHandler(env, env.notFound)
	return ep
}

func initIndexEndPoint(env *environment) endPoint {
	ep := newEndPoint("/")
	getRouting := newRoleBasedEndPointHandler(env, "navbar", "page.layout", "index")
	getRouting.defaultHandler = env.emptyProcessor
	ep.routing[http.MethodGet] = getRouting
	return ep
}

func initLoginEndPoint(env *environment) endPoint {
	ep := newEndPoint("/login")
	ep.routing[http.MethodGet] = newAuthEndPointHandler(
		env, env.login, "login.layout", "login")
	return ep
}

func initLogoutEndPoint(env *environment) endPoint {
	ep := newEndPoint("/logout")
	ep.routing[http.MethodGet] = newAuthEndPointHandler(env, env.logout)
	return ep
}

func initAuthEndPoint(env *environment) endPoint {
	ep := newEndPoint("/authenticate")
	ep.routing[http.MethodPost] = newAuthEndPointHandler(env, env.authenticate)
	return ep
}

func initCreateAccountEndPoint(env *environment) endPoint {
	ep := newEndPoint("/create_account")
	getRouting := newRoleBasedEndPointHandler(env, "login.layout", "create_account")
	getRouting.processors[models.AdminRole] = env.emptyProcessor
	ep.routing[http.MethodGet] = getRouting
	postRouting := newRoleBasedEndPointHandler(env)
	postRouting.processors[models.AdminRole] = env.createAccountPost
	ep.routing[http.MethodPost] = postRouting
	return ep
}

func initUsersEndPoint(env *environment) endPoint {
	ep := newEndPoint("/users")
	getRouting := newRoleBasedEndPointHandler(
		env, "navbar", "page.layout", "users")
	getRouting.processors[models.AdminRole] = env.users
	ep.routing[http.MethodGet] = getRouting
	return ep
}

func initProjectsEndPoint(env *environment) endPoint {
	ep := newEndPoint("/projects")
	getRouting := newRoleBasedEndPointHandler(
		env, "navbar", "page.layout", "projects")
	getRouting.processors[models.AdminRole] = env.emptyProcessor
	getRouting.processors[models.UserRole] = env.emptyProcessor
	ep.routing[http.MethodGet] = getRouting
	return ep
}

func initEchoProjectEndPoint(env *environment) endPoint {
	ep := newEndPoint("/projects/echo")
	getRouting := newRoleBasedEndPointHandler(
		env, "navbar", "page.layout", "echobot_chats")
	getRouting.processors[models.AdminRole] = env.echoProject
	getRouting.processors[models.UserRole] = env.echoProject
	ep.routing[http.MethodGet] = getRouting
	return ep
}
