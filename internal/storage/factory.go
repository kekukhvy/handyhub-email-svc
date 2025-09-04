package storage

import (
	"fmt"
	"handyhub-email-svc/internal/config"
	"handyhub-email-svc/internal/database"

	"github.com/sirupsen/logrus"
)

func NewEmailStorage(cfg *config.Configuration, mongodb *database.MongoDB) (EmailStorage, error) {
	switch cfg.Storage.Type {
	case "console":
		logrus.Info("Using Console Storage")
		return NewConsoleStorage(), nil
	case "file":
		logrus.Info("Using File Storage")
		return NewFileStorage(cfg.Storage.File)
	case "database":
		logrus.Info("Using Database Storage")
		return NewDatabaseStorage(mongodb, cfg.Database.EmailCollection)
	default:
		logrus.Warnf("Unknown storage type '%s', defaulting to Console Storage", cfg.Storage.Type)
		return nil, fmt.Errorf("unknown storage type: %s", cfg.Storage.Type)
	}
}
