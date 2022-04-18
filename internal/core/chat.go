package core

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `json:"id" bson:"id"`
	ChatId    primitive.ObjectID `json:"chat_id" bson:"chat_id"`
	Sender    primitive.ObjectID `json:"sender" bson:"sender"`
	Text      string             `json:"text" bson:"text"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type Chat struct {
	ID    primitive.ObjectID `json:"id" bson:"id"`
	FUser primitive.ObjectID `json:"f_user" bson:"f_user"`
	SUser primitive.ObjectID `json:"s_user" bson:"s_user"`
}
