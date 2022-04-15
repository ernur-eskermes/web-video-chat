package service

import (
	"context"
	"time"

	"github.com/ernur-eskermes/web-video-chat/pkg/room"
	"github.com/markbates/goth"
	"gopkg.in/olahol/melody.v1"

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
	AuthProvider(ctx context.Context, user goth.User) (Tokens, error)
	CreateSubscription(ctx context.Context, subscriberId, userId primitive.ObjectID) error
}

type Rooms interface {
	Create(ctx context.Context, input RoomCreateInput) (primitive.ObjectID, string, error)
}

type CreateMessageInput struct {
	ChatId  primitive.ObjectID
	UserId  primitive.ObjectID
	Message string
}

type Chats interface {
	GetMessages(ctx context.Context, chatId primitive.ObjectID) ([]domain.Message, error)
	CreateMessage(ctx context.Context, input CreateMessageInput) error
}

type Services struct {
	Users Users
	Rooms Rooms
	Chats Chats
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
	Websocket       *melody.Melody
}

func NewServices(deps Deps) *Services {
	usersService := NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager,
		deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.Domain)
	roomsService := NewRoomsService(deps.Repos.Rooms, deps.Room)
	chatsService := NewChatsService(deps.Repos.Chats, deps.Websocket)

	return &Services{
		Users: usersService,
		Rooms: roomsService,
		Chats: chatsService,
	}
}
