package echobot

import (
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/ShyBearStudio/tbot-admindashboard/logger"
)

const (
	configEnvVarName string = "echobotconfig"
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
	if len(configFileName) == 0 {
		return loadConfigFromEnv(configEnvVarName)
	}
	return loadConfigFromFile(configFileName)
}

func loadConfigFromFile(configFileName string) error {
	if len(configFileName) == 0 {
		return fmt.Errorf("Name of configuration file cannot be an empty string.")
	}
	file, err := os.Open(configFileName)
	if err != nil {
		return err
	}
	jsonDecoder := json.NewDecoder(file)
	err = jsonDecoder.Decode(&config)
	if err != nil {
		return err
	}
	return nil

}

func loadConfigFromEnv(envVarName string) error {
	//envVarValue := `{    "Address": "0.0.0.0:8080",    "StaticResources": "public",    "Db": {       "Driver": "postgres",       "ConnectionString": "postgres://leabqlqclrgvcs:546b03954b61761cc2e3b2c0dbbfae83936db2fca407af6495ef41e2ada899aa@ec2-54-163-246-165.compute-1.amazonaws.com:5432/d7g3dt1edh1ta2"    },    "Log": {       "Dir": "logs/admindashboard"    },    "tbots": {       "echobot": "configs/echobotconfig.json"    } }`
	envVarValue := os.Getenv(envVarName)
	byt := []byte(envVarValue)
	if err := json.Unmarshal(byt, &config); err != nil {
		panic(err)
	}
	return nil
}
