package repository

import (
	"context"
	"github.com/ernur-eskermes/web-video-chat/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatsRepo struct {
	db *mongo.Collection
}

func NewChatsRepo(db *mongo.Database) *ChatsRepo {
	return &ChatsRepo{
		db: db.Collection(chatsCollection),
	}
}

func (r *ChatsRepo) GetChatMessages(ctx context.Context, id primitive.ObjectID) ([]domain.Message, error) {
	var messages []domain.Message

	cur, err := r.db.Find(ctx, bson.M{"chat_id": id})
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *ChatsRepo) CreateMessage(ctx context.Context, message domain.Message) error {
	_, err := r.db.InsertOne(ctx, message)
	return err
}
