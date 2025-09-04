package storage

import (
	"handyhub-email-svc/internal/models"

	"github.com/sirupsen/logrus"
)

type ConsoleStorage struct{}

var log = logrus.StandardLogger()

func NewConsoleStorage() *ConsoleStorage {
	return &ConsoleStorage{}
}

func (cs *ConsoleStorage) Store(emailLog *models.EmailLog) error {
	logrus.WithFields(logrus.Fields{
		"to":       emailLog.To,
		"subject":  emailLog.Subject,
		"status":   emailLog.Status,
		"provider": emailLog.Provider,
		"attempts": emailLog.Attempts,
	}).Info("Email log entry")

	return nil
}

func (cs *ConsoleStorage) Close() error {
	log.Info("Console storage closed")
	return nil
}
