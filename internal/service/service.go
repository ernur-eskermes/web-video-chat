package service

import (
	"context"
	"github.com/ernur-eskermes/web-video-chat/pkg/room"
	"time"

	"github.com/ernur-eskermes/web-video-chat/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ernur-eskermes/web-video-chat/internal/repository"
	"github.com/ernur-eskermes/web-video-chat/pkg/auth"
	"github.com/ernur-eskermes/web-video-chat/pkg/hash"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type UserSignUpInput struct {
	Username        string
	Password        string
	ConfirmPassword string
}

type UserSignInInput struct {
	Username string
	Password string
}

type RoomCreateInput struct {
	UserId primitive.ObjectID
	Name   string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Users interface {
	SignUp(ctx context.Context, input UserSignUpInput) error
	SignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.User, error)
}

type Rooms interface {
	Create(ctx context.Context, input RoomCreateInput) (primitive.ObjectID, string, error)
}

type Services struct {
	Users Users
	Rooms Rooms
}

type Deps struct {
	Repos           *repository.Repositories
	Hasher          hash.PasswordHasher
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Environment     string
	Domain          string
	Room            room.Room
}

func NewServices(deps Deps) *Services {
	usersService := NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager,
		deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.Domain)
	roomsService := NewRoomsService(deps.Repos.Rooms, deps.Room)

	return &Services{
		Users: usersService,
		Rooms: roomsService,
	}
}
