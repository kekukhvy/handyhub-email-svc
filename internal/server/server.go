package server

import (
	"context"
	"errors"
	"handyhub-email-svc/internal/config"
	"handyhub-email-svc/internal/database"
	"handyhub-email-svc/internal/queue"
	"handyhub-email-svc/internal/smtp"
	"handyhub-email-svc/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

type Server struct {
	httpServer     *http.Server
	config         *config.Configuration
	mongodb        *database.MongoDB
	emailStorage   storage.EmailStorage
	rabbitMQ       *queue.RabbitMQ
	emailProcessor *queue.EmailProcessor
	smtpProvider   smtp.SMTPProvider
}

func New(cfg *config.Configuration) *Server {
	return &Server{
		config: cfg,
	}
}

func (s *Server) Start() error {
	if err := s.initMongoDB(); err != nil {
		return err
	}
	if err := s.initEmailStorage(); err != nil {
		return err
	}
	if err := s.initSMTPProvider(); err != nil {
		return err
	}
	if err := s.initRabbitMQ(); err != nil {
		return err
	}
	s.emailProcessor = queue.NewProcessor(s.emailStorage, s.smtpProvider)
	go s.startMessageConsumer()

	if err := s.setupHTTPServer(); err != nil {
		return err
	}

	go func() {
		log.Infof("Server starting on port %s", s.config.Server.Port)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Fatalf("Could not listen on %s: %v\n", s.config.Server.Port, err)
		}
	}()

	s.waitForShutdown()
	return nil
}

func (s *Server) initMongoDB() error {
	if s.config.Storage.Type == "database" {
		mongodb, err := database.NewMongoDB(*s.config)
		if err != nil {
			log.WithError(err).Fatal("Failed to initialize MongoDB")
			return err
		}
		s.mongodb = mongodb
	} else {
		log.Info("Skipping  MongoDB connection")
	}
	return nil
}

func (s *Server) initEmailStorage() error {
	emailStorage, err := storage.NewEmailStorage(s.config, s.mongodb)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize Email Storage")
		return err
	}
	s.emailStorage = emailStorage
	return nil
}

func (s *Server) initSMTPProvider() error {
	log.WithField("provider", s.config.SMTP.Provider).Info("Initializing SMTP Provider...")
	smtpProvider, err := smtp.NewSMTPProvider(s.config.SMTP)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize SMTP Provider")
		return err
	}
	s.smtpProvider = smtpProvider
	return nil
}

func (s *Server) initRabbitMQ() error {
	rabbitmq, err := queue.NewRabbitMQ(&s.config.Queue.RabbitMQ)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize RabbitMQ")
		return err
	}
	s.rabbitMQ = rabbitmq
	if err := rabbitmq.SetupQueue(); err != nil {
		log.WithError(err).Fatal("Failed to setup RabbitMQ queue")
		return err
	}
	return nil
}

func (s *Server) setupHTTPServer() error {
	gin.SetMode(s.config.Server.Mode)
	router := gin.Default()
	SetupRoutes(router, s.config, s.emailStorage)
	s.httpServer = &http.Server{
		Addr:         s.config.Server.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.config.Server.IdleTimeout) * time.Second,
	}
	log.Info("Initializing server...")
	return nil
}

func (s *Server) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.WithField("signal", sig).Info("Shutting down server...")

	s.Shutdown()
}

func (s *Server) startMessageConsumer() {
	log.Info("Starting message consumer...")
	messages, err := s.rabbitMQ.ConsumeMessages()

	if err != nil {
		log.WithError(err).Fatal("Failed to start consuming messages")
		return
	}

	for msg := range messages {
		log.Info("Received a new message")

		queueMessage, err := s.rabbitMQ.ParseMessage(msg.Body)
		if err != nil {
			log.WithError(err).Error("Failed to parse message, rejecting...")
			msg.Nack(false, false)
			continue
		}

		if err := s.emailProcessor.ProcessMessage(queueMessage); err != nil {
			log.WithError(err).Error("Failed to process message, rejecting...")
			msg.Nack(false, true)
			continue
		}

		msg.Ack(false)
		log.Info("Message processed and acknowledged")
	}

	log.Info("Message consumer stopped")
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if s.rabbitMQ != nil {
		if err := s.rabbitMQ.Close(); err != nil {
			log.WithError(err).Error("Error closing RabbitMQ connection")
		} else {
			log.Info("RabbitMQ connection closed")
		}
	}

	if s.emailStorage != nil {
		if err := s.emailStorage.Close(); err != nil {
			log.WithError(err).Error("Error closing email storage")
		} else {
			log.Info("Email storage closed")
		}
	}

	if s.mongodb != nil {
		if err := s.mongodb.Disconnect(ctx); err != nil {
			log.WithError(err).Error("Error disconnecting MongoDB")
		} else {
			log.Info("MongoDB disconnected")
		}
	}

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server forced to shutdown")
	}

	log.Info("Server gracefully stopped")
}
