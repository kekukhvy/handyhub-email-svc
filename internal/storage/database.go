package storage

import (
	"context"
	"handyhub-email-svc/internal/database"
	"handyhub-email-svc/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type DatabaseStorage struct {
	mongodb    *database.MongoDB
	collection *mongo.Collection
}

func NewDatabaseStorage(mongodb *database.MongoDB, collectionName string) (*DatabaseStorage, error) {
	collection := mongodb.Database.Collection(collectionName)

	log.Info("Database storage initialized")

	return &DatabaseStorage{
		mongodb:    mongodb,
		collection: collection,
	}, nil
}

func (ds *DatabaseStorage) Store(emailLog *models.EmailLog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := ds.collection.InsertOne(ctx, emailLog)
	if err != nil {
		log.WithError(err).Error("Failed to store email log in database")
		return err
	}

	log.Info("Email log entry stored in database")

	return nil
}

func (ds *DatabaseStorage) Close() error {
	log.Info("Database storage closed")
	return nil
}
