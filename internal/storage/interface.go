package storage

import "handyhub-email-svc/internal/models"

type EmailStorage interface {
	Store(emailLog *models.EmailLog) error
	Close() error
}
