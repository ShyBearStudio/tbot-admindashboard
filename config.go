package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Address         string
	StaticResources string
	DbConfig        struct {
		Driver           string
		ConnectionString string
	}
}

var config Config

func loadConfig(configFileName string) error {
	_ = "breakpoint"
	file, err := os.Open(configFileName)
	if err != nil {
		log.Fatalln("Cannot open configuration file", err)
		return err
	}
	jsonDecoder := json.NewDecoder(file)
	config = Config{}
	err = jsonDecoder.Decode(&config)
	if err != nil {
		log.Fatalf("Cannot get configuration from file '%s'...\n", configFileName)
		return err
	}

	return nil
}
