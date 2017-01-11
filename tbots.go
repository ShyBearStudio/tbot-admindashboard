package main

import (
	"fmt"

	"github.com/ShyBearStudio/tbot-admindashboard/logger"
	"github.com/ShyBearStudio/tbot-admindashboard/projects/echobot"
	"github.com/ShyBearStudio/tbot-admindashboard/projects/tbot"
)

type tbotId string

const (
	undefinedBotId tbotId = ""
	echoBotId      tbotId = "echobot"
)

func newTBots(tbotConfigs map[string]string) (map[tbotId]tbot.TBot, error) {
	tbots := make(map[tbotId]tbot.TBot)
	for tbotName, configFileName := range tbotConfigs {
		tbotId, err := toTBotId(tbotName)
		if err != nil {
			return nil, err
		}
		tbot, err := newTBot(tbotId, configFileName)
		if err != nil {
			logger.Errorf("Cannot initialized bot (Name: '%s' | Config: '%s'): %v", tbotName, configFileName, err)
			return nil, err
		}
		tbots[tbotId] = tbot
	}

	return tbots, nil
}

func toTBotId(s string) (id tbotId, err error) {
	if s == string(echoBotId) {
		return echoBotId, nil
	}

	return undefinedBotId, fmt.Errorf("Unknown tbot id value: '%s'", s)
}

func newTBot(id tbotId, configFileName string) (tbot.TBot, error) {
	switch id {
	case echoBotId:
		return echobot.New(configFileName)
	}
	return nil, fmt.Errorf("There is no way to initialize tbot with id '%s'", id)
}

func (env *environment) startTBots() error {
	for id, tbot := range env.tbots {
		err := tbot.Start()
		if err != nil {
			logger.Errorf("Cannot start tbot (id: '%s'): '%v'", id, err)
			return err
		}
	}

	return nil
}
