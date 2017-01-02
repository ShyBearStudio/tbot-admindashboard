package main

import (
	"flag"
	"net/http"

	"github.com/ShyBearStudio/tbot-admindashboard/data"
	"github.com/ShyBearStudio/tbot-admindashboard/logger"
	"github.com/gorilla/mux"
)

var cmdLineArgs commandLineAgruments

func init() {
	cmdLineArgs = commandLineAgruments{}
	cmdLineArgs.help = flag.Bool("help", false, "get help")
	cmdLineArgs.config = flag.String("config", "config.json", "config file in JSON format")

}

const (
	staticResourcesPrefix string = "/static/"
)

func main() {
	_ = "breakpoint"
	flag.Parse()
	if err := loadConfig(*cmdLineArgs.config); err != nil {
		return
	}
	if err := logger.InitLog(config.Log.Dir); err != nil {
		return
	}
	if err := data.InitDb(config.Db.Driver, config.Db.ConnectionString); err != nil {
		return
	}

	r := registerEndPoints()
	mux := http.NewServeMux()
	// Add static resources
	staticFiles := http.FileServer(http.Dir(config.StaticResources))
	mux.Handle(
		staticResourcesPrefix, http.StripPrefix(staticResourcesPrefix, staticFiles))
	mux.Handle("/", r)

	server := &http.Server{
		Addr:           config.Address,
		Handler:        mux,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}

func registerEndPoints() *mux.Router {
	r := mux.NewRouter()
	for _, ep := range endPoints {
		for method, handler := range ep.routing {
			r.Handle(ep.pattern, handler).Methods(method)
		}
	}

	r.NotFoundHandler = http.HandlerFunc(notFount)

	return r
}
