package common

import (
	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

func init() {
	// do something here to set environment depending on an environment variable
	// or command-line flag
	if Config.App.Environment == "production" {
		logger.SetFormatter(&log.JSONFormatter{})
	} else {
		// The TextFormatter is default, you don't actually have to do this.
		logger.SetFormatter(&log.TextFormatter{})
	}
}
