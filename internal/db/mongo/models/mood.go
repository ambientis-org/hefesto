package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Mood struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Value     int                `bson:"value"`
}

func NewMood(value int) Mood {
	return Mood{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Value:     value,
	}
}
