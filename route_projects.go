package main

import (
	"net/http"

	"github.com/ShyBearStudio/tbot-admindashboard/logger"
	"github.com/ShyBearStudio/tbot-admindashboard/projects/echobot/data"
)

func projects(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, nil, "page.layout", "projects")
}

func echoProject(w http.ResponseWriter, r *http.Request) {
	//echoBot := tbots[echoBotId].(*echobot.EchoBot)
	chats, err := data.Chats()
	if err != nil {
		logger.Errorf("Cannot get chats: %v", err)
	}
	generateHTML(w, &chats, "page.layout", "echobot_chats")
}
