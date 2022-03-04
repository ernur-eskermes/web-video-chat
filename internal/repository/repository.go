package repository

import (
	"context"

	"github.com/ernur-eskermes/web-video-chat/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Users interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	SetSession(ctx context.Context, userID primitive.ObjectID, session domain.Session) error
	GetById(ctx context.Context, id primitive.ObjectID) (domain.User, error)
}

type Rooms interface {
	Create(ctx context.Context, room domain.Room) (primitive.ObjectID, error)
}

type Repositories struct {
	Users Users
	Rooms Rooms
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Users: NewUsersRepo(db),
		Rooms: NewRoomsRepo(db),
	}
}
