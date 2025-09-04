package smtp

import (
	"fmt"
	"handyhub-email-svc/internal/config"
)

func NewSMTPProvider(cfg config.SMTPConfig) (SMTPProvider, error) {
	switch cfg.Provider {
	case "gmail":
		if cfg.Gmail.Username == "" || cfg.Gmail.Password == "" {
			return nil, fmt.Errorf("gmail provider requires username and password")
		}
		return NewGmailProvider(cfg.Gmail, cfg.DefaultFrom), nil
	case "sendgrid":
		if cfg.SendGrid.ApiKey == "" || cfg.SendGrid.Url == "" {
			return nil, fmt.Errorf("sendgrid provider requires api key and url")
		}
		return NewSendGridProvider(cfg.SendGrid, cfg.DefaultFrom), nil

	case "mailhog":
		if cfg.MailHog.Host == "" || cfg.MailHog.Port == 0 {
			return nil, fmt.Errorf("mailhog provider requires host and port")
		}
		return NewMailHogProvider(cfg.MailHog, cfg.DefaultFrom), nil
	default:
		return nil, fmt.Errorf("unsupported SMTP provider: %s", cfg.Provider)
	}
}
