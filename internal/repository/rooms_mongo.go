package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/ernur-eskermes/web-video-chat/internal/core"
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

func (r *RoomsRepo) Create(ctx context.Context, room core.Room) (primitive.ObjectID, error) {
	res, err := r.db.InsertOne(ctx, room)

	return res.InsertedID.(primitive.ObjectID), err
}

func (r *RoomsRepo) GetList(ctx context.Context, roomVisibility bool) ([]core.Room, error) {
	var rooms []core.Room

	cur, err := r.db.Find(ctx, bson.M{"visibility": roomVisibility})
	if err != nil {
		return nil, err
	}

	if err = cur.All(ctx, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *RoomsRepo) GetById(ctx context.Context, roomID primitive.ObjectID) (core.Room, error) {
	var room core.Room

	if err := r.db.FindOne(ctx, bson.M{"_id": roomID}).Decode(&room); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.Room{}, core.ErrRoomNotFound
		}

		return core.Room{}, err
	}

	return room, nil
}
