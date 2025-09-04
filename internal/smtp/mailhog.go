package smtp

import (
	"fmt"
	"handyhub-email-svc/internal/config"
	"handyhub-email-svc/internal/models"

	"gopkg.in/gomail.v2"
)

type MailHogProvider struct {
	config config.MailHogConfig
	from   string
}

func NewMailHogProvider(cfg config.MailHogConfig, from string) *MailHogProvider {
	return &MailHogProvider{
		config: cfg,
		from:   from,
	}
}

func (m *MailHogProvider) SendEmail(email *models.EmailMessage) error {
	if len(email.To) == 0 {
		return fmt.Errorf("no recipients specified")
	}
	if email.BodyHTML == "" && email.BodyText == "" {
		return fmt.Errorf("email body is required")
	}

	msg := gomail.NewMessage()
	m.setHeaders(msg, email)

	if email.BodyHTML != "" {
		msg.SetBody("text/html", email.BodyHTML)
		if email.BodyText != "" {
			msg.AddAlternative("text/plain", email.BodyText)
		}
	} else {
		msg.SetBody("text/plain", email.BodyText)
	}

	d := gomail.NewDialer(m.config.Host, m.config.Port, "", "")
	d.TLSConfig = nil

	if err := d.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email via MailHog: %w", err)
	}
	return nil
}

func (m *MailHogProvider) setHeaders(msg *gomail.Message, email *models.EmailMessage) {
	fromEmail := email.From
	if fromEmail == "" {
		fromEmail = m.from
	}
	msg.SetHeader("From", fromEmail)
	msg.SetHeader("To", email.To...)
	msg.SetHeader("Subject", email.Subject)
	msg.SetHeader("X-Mailer", "HandyHub Email Service")
	msg.SetHeader("X-Environment", "development")
}

func (m *MailHogProvider) GetProviderName() string {
	return "mailhog"
}
