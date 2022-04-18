package repository

import (
	"context"

	"github.com/ernur-eskermes/web-video-chat/internal/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Users interface {
	Create(ctx context.Context, user *core.User) error
	GetByCredentials(ctx context.Context, email, password, provider string) (core.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (core.User, error)
	SetSession(ctx context.Context, userID primitive.ObjectID, session core.Session) error
	GetById(ctx context.Context, id primitive.ObjectID) (core.User, error)
	CreateSubscription(ctx context.Context, userId, subscription primitive.ObjectID) error
}

type Rooms interface {
	Create(ctx context.Context, room core.Room) (primitive.ObjectID, error)
	GetList(ctx context.Context, roomVisibility bool) ([]core.Room, error)
	GetById(ctx context.Context, roomID primitive.ObjectID) (core.Room, error)
}

type Chats interface {
	GetChatMessages(ctx context.Context, id primitive.ObjectID) ([]core.Message, error)
	CreateMessage(ctx context.Context, message core.Message) error
}

type Repositories struct {
	Users Users
	Rooms Rooms
	Chats Chats
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Users: NewUsersRepo(db),
		Rooms: NewRoomsRepo(db),
		Chats: NewChatsRepo(db),
	}
}
