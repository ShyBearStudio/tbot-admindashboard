package main

import (
	"net/http"

	"github.com/ShyBearStudio/tbot-admindashboard/logger"
	"github.com/ShyBearStudio/tbot-admindashboard/projects/echobot"
)

func (env *environment) echoProject(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	echoBot := env.tbots[echoBotId].(*echobot.EchoBot)
	chats, err := echoBot.Chats()
	if err != nil {
		logger.Errorf("Cannot get chats: %v", err)
	}
	return chats, err
}
