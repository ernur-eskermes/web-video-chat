package repository

import (
	"context"
	"github.com/ernur-eskermes/web-video-chat/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomsRepo struct {
	db *mongo.Collection
}

func NewRoomsRepo(db *mongo.Database) *RoomsRepo {
	return &RoomsRepo{
		db: db.Collection(roomsCollection),
	}
}

func (r *RoomsRepo) Create(ctx context.Context, room domain.Room) (primitive.ObjectID, error) {
	res, err := r.db.InsertOne(ctx, room)

	return res.InsertedID.(primitive.ObjectID), err
}
