package queue

import (
	"handyhub-email-svc/internal/models"
	"handyhub-email-svc/internal/smtp"
	"handyhub-email-svc/internal/storage"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmailProcessor struct {
	emailStorage storage.EmailStorage
	smtpProvider smtp.SMTPProvider
}

func NewProcessor(emailStorage storage.EmailStorage, smtpProvider smtp.SMTPProvider) *EmailProcessor {
	return &EmailProcessor{
		emailStorage: emailStorage,
		smtpProvider: smtpProvider,
	}
}

func (p *EmailProcessor) ProcessMessage(message *models.QueueMessage) error {
	log.WithFields(logrus.Fields{
		"to":       message.Email.To,
		"subject":  message.Email.Subject,
		"provider": p.smtpProvider.GetProviderName(),
	}).Info("Processing email message")

	var emailLog *models.EmailLog
	emailLog = &models.EmailLog{
		ID:       primitive.NewObjectID(),
		To:       message.Email.To,
		Subject:  message.Email.Subject,
		Provider: p.smtpProvider.GetProviderName(),
		Attempts: 1,
		SentAt:   time.Now(),
	}

	err := p.smtpProvider.SendEmail(&message.Email)
	if err != nil {
		log.WithError(err).Error("Failed to send email")
		emailLog.Status = "failed"
		emailLog.ErrorMsg = err.Error()

	} else {
		log.Info("Email sent successfully")
		emailLog.Status = "success"
	}

	if err := p.emailStorage.Store(emailLog); err != nil {
		log.WithError(err).Error("Failed to store email log")
		return err
	}

	log.Info("Email processed and logged successfully")
	return nil
}
