package main

import (
	"handyhub-email-svc/internal/config"
	"handyhub-email-svc/internal/logger"
	"handyhub-email-svc/internal/server"

	"github.com/sirupsen/logrus"
)

var log = *logrus.StandardLogger()

func main() {

	cfg := config.Load()
	logger.Init(cfg.Logs.Level, cfg.Logs.Path)

	log.Infof("Application Name %s is starting....", cfg.App.Name)
	srv := server.New(cfg)
	if err := srv.Start(); err != nil {
		log.WithError(err).Fatalf("Error starting server: %v", err)
	}
}
