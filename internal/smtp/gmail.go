package smtp

import (
	"fmt"
	"handyhub-email-svc/internal/config"
	"handyhub-email-svc/internal/models"

	"gopkg.in/gomail.v2"
)

type GmailProvider struct {
	config config.GmailConfig
	from   string
}

func NewGmailProvider(cfg config.GmailConfig, from string) *GmailProvider {
	return &GmailProvider{
		config: cfg,
		from:   from,
	}
}

func (g *GmailProvider) SendEmail(email *models.EmailMessage) error {
	m := gomail.NewMessage()

	fromEmail := email.From
	if fromEmail == "" {
		fromEmail = g.from
	}
	m.SetHeader("From", fromEmail)

	if len(email.To) == 0 {
		return fmt.Errorf("no recipients specified")
	}
	m.SetHeader("To", email.To...)
	m.SetHeader("Subject", email.Subject)

	if email.BodyHTML != "" {
		m.SetBody("text/html", email.BodyHTML)
		if email.BodyText != "" {
			m.AddAlternative("text/plain", email.BodyText)
		}
	} else if email.BodyText != "" {
		m.SetBody("text/plain", email.BodyText)
	} else {
		return fmt.Errorf("email body is required")
	}

	d := gomail.NewDialer(g.config.Host, g.config.Port, g.config.Username, g.config.Password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email via Gmail: %w", err)
	}

	return nil
}

func (g *GmailProvider) GetProviderName() string {
	return "gmail"
}
