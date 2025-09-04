package smtp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"handyhub-email-svc/internal/config"
	"handyhub-email-svc/internal/models"
	"net/http"
)

type SendGridProvider struct {
	config config.SendGridConfig
	from   string
}

type sendGridMessage struct {
	Personalizations []sendGridPersonalization `json:"personalizations"`
	From             sendGridEmail             `json:"from"`
	Subject          string                    `json:"subject"`
	Content          []sendGridContent         `json:"content"`
}
type sendGridPersonalization struct {
	To []sendGridEmail `json:"to"`
}

type sendGridEmail struct {
	Email string `json:"email"`
}

type sendGridContent struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func NewSendGridProvider(cfg config.SendGridConfig, from string) *SendGridProvider {
	return &SendGridProvider{
		config: cfg,
		from:   from,
	}
}

func (s *SendGridProvider) SendEmail(email *models.EmailMessage) error {
	if len(email.To) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	to := s.buildRecipients(email.To)
	content, err := s.buildContent(email)
	if err != nil {
		return err
	}

	fromEmail := email.From
	if fromEmail == "" {
		fromEmail = s.from
	}

	message := s.buildMessage(to, fromEmail, email.Subject, content)
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal SendGrid message: %w", err)
	}

	req, err := s.createRequest(jsonData)
	if err != nil {
		return err
	}

	s.setHeaders(req)
	resp, err := s.sendRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return s.handleResponse(resp)
}

func (s *SendGridProvider) buildRecipients(recipients []string) []sendGridEmail {
	to := make([]sendGridEmail, 0, len(recipients))
	for _, recipient := range recipients {
		to = append(to, sendGridEmail{Email: recipient})
	}
	return to
}

func (s *SendGridProvider) buildContent(email *models.EmailMessage) ([]sendGridContent, error) {
	var content []sendGridContent
	if email.BodyText != "" {
		content = append(content, sendGridContent{
			Type:  "text/plain",
			Value: email.BodyText,
		})
	}
	if email.BodyHTML != "" {
		content = append(content, sendGridContent{
			Type:  "text/html",
			Value: email.BodyHTML,
		})
	}
	if len(content) == 0 {
		return nil, fmt.Errorf("email body is required")
	}
	return content, nil
}

func (s *SendGridProvider) buildMessage(to []sendGridEmail, from, subject string, content []sendGridContent) sendGridMessage {
	return sendGridMessage{
		Personalizations: []sendGridPersonalization{{To: to}},
		From:             sendGridEmail{Email: from},
		Subject:          subject,
		Content:          content,
	}
}

func (s *SendGridProvider) createRequest(jsonData []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", s.config.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	return req, nil
}

func (s *SendGridProvider) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+s.config.ApiKey)
	req.Header.Set("Content-Type", "application/json")
}

func (s *SendGridProvider) sendRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to SendGrid: %w", err)
	}
	return resp, nil
}

func (s *SendGridProvider) handleResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("SendGrid API returned status %d", resp.StatusCode)
	}
	return nil
}

func (s *SendGridProvider) GetProviderName() string {
	return "sendgrid"
}
