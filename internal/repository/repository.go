package repository

import (
	"context"

	"github.com/ernur-eskermes/web-video-chat/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Users interface {
	Create(ctx context.Context, user *domain.User) error
	GetByCredentials(ctx context.Context, email, password, provider string) (domain.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	SetSession(ctx context.Context, userID primitive.ObjectID, session domain.Session) error
	GetById(ctx context.Context, id primitive.ObjectID) (domain.User, error)
	CreateSubscription(ctx context.Context, userId, subscription primitive.ObjectID) error
}

type Rooms interface {
	Create(ctx context.Context, room domain.Room) (primitive.ObjectID, error)
}

type Chats interface {
	GetChatMessages(ctx context.Context, id primitive.ObjectID) ([]domain.Message, error)
	CreateMessage(ctx context.Context, message domain.Message) error
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
