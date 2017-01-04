package main

import (
	"fmt"

	"github.com/ShyBearStudio/tbot-admindashboard/logger"
	"github.com/ShyBearStudio/tbot-admindashboard/projects/echobot"
)

type tbotId string

const (
	undefinedBotId tbotId = ""
	echoBotId      tbotId = "echobot"
)

var tbots = make(map[tbotId]echobot.TBot)

func initTBots(tbotConfigs map[string]string) error {
	for tbotName, configFileName := range tbotConfigs {
		tbotId, err := toTBotId(tbotName)
		if err != nil {
			return err
		}
		tbot, err := initTBot(tbotId, configFileName)
		if err != nil {
			logger.Errorf("Cannot initialized bot (Name: '%s' | Config: '%s'): %v", tbotName, configFileName, err)
			return err
		}
		tbots[tbotId] = tbot
	}

	return nil
}

func toTBotId(s string) (id tbotId, err error) {
	if s == string(echoBotId) {
		return echoBotId, nil
	}

	return undefinedBotId, fmt.Errorf("Unknown tbot id value: '%s'", s)
}

func initTBot(id tbotId, configFileName string) (echobot.TBot, error) {
	switch id {
	case echoBotId:
		return echobot.New(configFileName)
	}
	return nil, fmt.Errorf("There is no way to initialize tbot with id '%s'", id)
}

func startTBots() error {
	for id, tbot := range tbots {
		err := tbot.Start()
		if err != nil {
			logger.Errorf("Cannot start tbot (id: '%s'): '%v'", id, err)
			return err
		}
	}

	return nil
}
