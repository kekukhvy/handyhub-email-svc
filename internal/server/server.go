package server

import (
	"context"
	"handyhub-email-svc/internal/config"
	"handyhub-email-svc/internal/database"
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
	httpServer   *http.Server
	config       *config.Configuration
	mongodb      *database.MongoDB
	emailStorage storage.EmailStorage
}

func New(cfg *config.Configuration) *Server {
	return &Server{
		config: cfg,
	}
}

func (s *Server) Start() error {

	var err error
	var mongodb *database.MongoDB

	if s.config.Storage.Type == "database" {
		mongodb, err = database.NewMongoDB(*s.config)
		if err != nil {
			log.WithError(err).Fatal("Failed to initialize MongoDB")
			return err
		}
		s.mongodb = mongodb
	} else {
		log.Info("Skipping  MongoDB connection")
	}

	emailStorage, err := storage.NewEmailStorage(s.config, mongodb)

	if err != nil {
		log.WithError(err).Fatal("Failed to initialize Email Storage")
		return err
	}
	s.emailStorage = emailStorage

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

	go func() {
		log.Infof("Server starting on port %s", s.config.Server.Port)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatalf("Could not listen on %s: %v\n", s.config.Server.Port, err)
		}
	}()

	s.waitForShutdown()

	return nil
}

func (s *Server) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.WithField("signal", sig).Info("Shutting down server...")

	s.Shutdown()
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
