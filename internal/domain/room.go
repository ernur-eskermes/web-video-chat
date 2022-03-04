package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name      string               `json:"name" bson:"name"`
	Author    primitive.ObjectID   `json:"author" bson:"author"`
	Members   []primitive.ObjectID `json:"members" bson:"members"`
	CreatedAt time.Time            `json:"created_at" bson:"created_at"`
}
