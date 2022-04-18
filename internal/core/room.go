package core

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID         primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name       string               `json:"name" bson:"name"`
	Author     primitive.ObjectID   `json:"author" bson:"author"`
	Visibility bool                 `json:"visibility" bson:"visibility"`
	Members    []primitive.ObjectID `json:"members" bson:"members"`
	Invites    []primitive.ObjectID `json:"invites" bson:"invites"`
	CreatedAt  time.Time            `json:"created_at" bson:"created_at"`
}

type RoomCreateInput struct {
	Name       string               `json:"name" binding:"required,min=8,max=64"`
	Visibility bool                 `json:"visibility" binding:"required"`
	Invites    []primitive.ObjectID `json:"invites"`
}
