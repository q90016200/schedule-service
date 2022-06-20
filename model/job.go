package model

import "time"

type Job struct {
	Name      string    `json:"name" bson:"name"`
	Method    string    `json:"method" bson:"method"`
	Path      string    `json:"path" bson:"path"`
	Cron      string    `json:"cron" bson:"cron"`
	Status    string    `json:"status" bson:"status"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
