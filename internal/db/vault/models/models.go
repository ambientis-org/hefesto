package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Mood struct {
	ID 		  primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	Value 	  int 		`bson:"value"`
}

func NewMood(value int) Mood {
	return Mood{
		ID: primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Value: value,
	}
}

type Journal struct {
	ID 		  primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	UserID	  uint		`bson:"user_id"`
	Username  string	`bson:"username"`
	Moods	  []Mood	`bson:"moods"`
}

func NewJournal(userID uint, username string) *Journal {
	return &Journal{
		ID: primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: userID,
		Username: username,
		Moods: []Mood{},
	}
}