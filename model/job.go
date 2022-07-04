package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Job struct {
	ID        primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name,omitempty"`
	Method    string    `json:"method" bson:"method,omitempty"`
	Path      string    `json:"path" bson:"path,omitempty"`
	Cron      string    `json:"cron" bson:"cron,omitempty"`
	Status    string    `json:"status" bson:"status,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
