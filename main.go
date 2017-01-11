package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/ShyBearStudio/tbot-admindashboard/configutils"
	"github.com/ShyBearStudio/tbot-admindashboard/logger"
	"github.com/ShyBearStudio/tbot-admindashboard/models"
	"github.com/gorilla/mux"
)

var cmdLineArgs commandLineAgruments
var config = configuration{}

type environment struct {
	config    configuration
	db        models.Datastore
	endPoints map[endPointId]endPoint
}

func init() {
	//cmdLineArgs = commandLineAgruments{}
	//cmdLineArgs.help = flag.Bool("help", false, "get help")
	//cmdLineArgs.config = flag.String("config", "", "config file in JSON format")
}

const (
	staticResourcesPrefix string = "/static/"
)

func main() {
	_ = "breakpoint"
	flag.Parse()
	env, err := newEnvironment()
	if err != nil {
		logger.Errorf("Cannot create environment: %v", err)
		return
	}

	if err := configutils.LoadConfig(*cmdLineArgs.config, configEnvVarName, &config); err != nil {
		return
	}
	if err := logger.InitLogger(config.Log.Dir); err != nil {
		return
	}
	//if err := data.InitDb(config.Db.Driver, config.Db.ConnectionString); err != nil {
	//	return
	//}
	if err := initTBots(config.Tbots); err != nil {
		return
	}
	if err := startTBots(); err != nil {
		return
	}

	r := registerEndPoints(env)
	mux := http.NewServeMux()
	// Add static resources
	staticFiles := http.FileServer(http.Dir(config.StaticResources))
	mux.Handle(
		staticResourcesPrefix, http.StripPrefix(staticResourcesPrefix, staticFiles))
	mux.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        mux,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}

func newEnvironment() (env *environment, err error) {
	cmdLineArgs = newCmdLineArgs()
	env = new(environment)
	if err = configutils.LoadConfig(*cmdLineArgs.config, configEnvVarName, &env.config); err != nil {
		logger.Errorf("Cannot load config: %v", err)
		return nil, err
	}
	if err = logger.InitLogger(env.config.Log.Dir); err != nil {
		logger.Errorf("Cannot initialize logger: %v", err)
		return nil, err
	}
	if env.db, err = models.NewDb(env.config.Db.Driver, env.config.Db.ConnectionString); err != nil {
		logger.Errorf("Cannot create database driver: %v", err)
		return nil, err
	}
	env.endPoints = newEndpoins(env)
	return
}

func newCmdLineArgs() commandLineAgruments {
	cmdLineArgs = commandLineAgruments{}
	cmdLineArgs.help = flag.Bool("help", false, "get help")
	cmdLineArgs.config = flag.String("config", "", "config file in JSON format")
	flag.Parse()
	return cmdLineArgs
}

func registerEndPoints(env *environment) *mux.Router {
	r := mux.NewRouter()
	for _, ep := range env.endPoints {
		for method, handler := range ep.routing {
			r.Handle(ep.pattern, handler).Methods(method)
		}
	}

	r.NotFoundHandler = http.HandlerFunc(env.notFount)

	return r
}
