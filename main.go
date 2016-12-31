package main

import (
	"flag"
	"net/http"

	"github.com/ShyBearStudio/tbot-admindashboard/data"
	"github.com/ShyBearStudio/tbot-admindashboard/logger"
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

	mux := http.NewServeMux()
	// Add static resources
	staticFiles := http.FileServer(http.Dir(config.StaticResources))
	mux.Handle(
		staticResourcesPrefix, http.StripPrefix(staticResourcesPrefix, staticFiles))
	// Configure routing
	AddEndPoint(mux, endPoints.index)
	AddEndPoint(mux, endPoints.login)

	server := &http.Server{
		Addr:           config.Address,
		Handler:        mux,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}

func AddEndPoint(mux *http.ServeMux, endPoint endPointDesc) {
	mux.Handle(endPoint.pattern, &endPoint)
}
