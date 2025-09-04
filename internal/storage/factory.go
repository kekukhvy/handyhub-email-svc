package storage

import (
	"fmt"
	"handyhub-email-svc/internal/config"

	"github.com/sirupsen/logrus"
)

func newEmailStorage(cfg *config.Configuration) (EmailStorage, error) {
	switch cfg.Storage.Type {
	case "console":
		logrus.Info("Using Console Storage")
		return NewConsoleStorage(), nil
	default:
		logrus.Warnf("Unknown storage type '%s', defaulting to Console Storage", cfg.Storage.Type)
		return nil, fmt.Errorf("unknown storage type: %s", cfg.Storage.Type)
	}
}
