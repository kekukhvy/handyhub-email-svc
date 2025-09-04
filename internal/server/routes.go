package server

import (
	"handyhub-email-svc/internal/config"
	"handyhub-email-svc/internal/models"
	"handyhub-email-svc/internal/storage"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var logger = logrus.StandardLogger()

func SetupRoutes(router *gin.Engine, cfg *config.Configuration, emailStorage storage.EmailStorage) {

	// Health endpoint
	router.GET("/health", func(c *gin.Context) {
		logrus.Info("Health check requested")
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "email-service",
		})
	})

	api := router.Group("/api/v1")
	{
		api.GET("/status", func(c *gin.Context) {
			logrus.Info("API status requested")
			c.JSON(200, gin.H{
				"api_version": "v1",
				"status":      "operational",
			})
		})

		// Тестовый endpoint для проверки хранения email логов
		api.POST("/test-email-log", func(c *gin.Context) {
			logger.Info("Test email log requested")

			// Создать тестовый email log
			testLog := &models.EmailLog{
				ID:       primitive.NewObjectID(),
				To:       []string{"test@example.com"},
				Subject:  "Test Email Log",
				Status:   "success",
				Provider: "test",
				Attempts: 1,
				SentAt:   time.Now(),
				ErrorMsg: "",
			}

			// Сохранить через выбранное хранилище
			if err := emailStorage.Store(testLog); err != nil {
				logger.WithError(err).Error("Failed to store test email log")
				c.JSON(500, gin.H{
					"error": "Failed to store email log",
				})
				return
			}

			logger.Info("Test email log stored successfully")
			c.JSON(200, gin.H{
				"message":      "Test email log stored",
				"storage_type": cfg.Storage.Type,
				"log":          testLog,
			})
		})
	}
}
