package service

import (
	"context"
	"github.com/ernur-eskermes/web-video-chat/internal/domain"
	"github.com/ernur-eskermes/web-video-chat/internal/repository"
	"github.com/ernur-eskermes/web-video-chat/pkg/room"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type RoomsService struct {
	repo repository.Rooms
	room room.Room
}

func NewRoomsService(repo repository.Rooms, room room.Room) *RoomsService {
	return &RoomsService{
		repo: repo,
		room: room,
	}
}

func (r *RoomsService) Create(ctx context.Context, input RoomCreateInput) (primitive.ObjectID, string, error) {
	roomId, err := r.repo.Create(ctx, domain.Room{
		Name:      input.Name,
		Author:    input.UserId,
		Members:   []primitive.ObjectID{input.UserId},
		CreatedAt: time.Now(),
	})
	if err != nil {
		return primitive.ObjectID{}, "", err
	}
	token, err := r.room.GetJoinToken(roomId.String(), input.UserId.String())
	if err != nil {
		return primitive.ObjectID{}, "", err
	}
	return roomId, token, nil
}
