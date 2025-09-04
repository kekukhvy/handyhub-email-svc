package server

import (
	"handyhub-email-svc/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SetupRoutes(router *gin.Engine, cfg *config.Configuration) {

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
	}
}
