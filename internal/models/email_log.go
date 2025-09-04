package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmailLog struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	To       []string           `json:"to" bson:"to"`
	Subject  string             `json:"subject" bson:"subject"`
	Status   string             `json:"status" bson:"status"`
	Provider string             `json:"provider" bson:"provider"`
	Attempts int                `json:"attempts" bson:"attempts"`
	SentAt   time.Time          `json:"sent_at" bson:"sent_at"`
	ErrorMsg string             `json:"error_msg,omitempty" bson:"error_msg,omitempty"`
}
