package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Username     string               `json:"username" bson:"username"`
	Password     string               `json:"password" bson:"password"`
	RegisteredAt time.Time            `json:"registeredAt" bson:"registeredAt"`
	LastVisitAt  time.Time            `json:"lastVisitAt" bson:"lastVisitAt"`
	Provider     string               `json:"provider" bson:"provider"`
	ACCSubs      []primitive.ObjectID `json:"acc_subs" bson:"acc_subs"`
	PNDSubs      []primitive.ObjectID `json:"pnd_subs" bson:"pnd_subs"`
}
