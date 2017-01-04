package echobot

import (
	"encoding/json"
	"os"

	"github.com/ShyBearStudio/tbot-admindashboard/logger"
)

type configuration struct {
	Token string
	Db    struct {
		Driver           string
		ConnectionString string
	}
	Log struct {
		Dir string
	}
}

var config configuration = configuration{}

func loadConfig(configFileName string) error {
	_ = "breakpoint"
	file, err := os.Open(configFileName)
	if err != nil {
		logger.Errorln("Cannot open configuration file", err)
		return err
	}
	jsonDecoder := json.NewDecoder(file)
	err = jsonDecoder.Decode(&config)
	if err != nil {
		logger.Errorln("Cannot get configuration from file", err)
		return err
	}
	return nil
}
