package smtp

import "handyhub-email-svc/internal/models"

type SMTPProvider interface {
	SendEmail(email *models.EmailMessage) error
	GetProviderName() string
}
