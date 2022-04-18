package service

import (
	"context"
	"time"

	"github.com/ernur-eskermes/web-video-chat/internal/core"
	"github.com/ernur-eskermes/web-video-chat/internal/repository"
	"github.com/ernur-eskermes/web-video-chat/pkg/room"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (s *RoomsService) Create(ctx context.Context, input core.RoomCreateInput, userID primitive.ObjectID) (primitive.ObjectID, string, error) {
	roomId, err := s.repo.Create(ctx, core.Room{
		Name:       input.Name,
		Visibility: input.Visibility,
		Author:     userID,
		Members:    []primitive.ObjectID{userID},
		CreatedAt:  time.Now(),
	})
	if err != nil {
		return primitive.ObjectID{}, "", err
	}

	token, err := s.room.GetJoinToken(roomId.String(), userID.String())
	if err != nil {
		return primitive.ObjectID{}, "", err
	}

	return roomId, token, nil
}

func (s *RoomsService) GetList(ctx context.Context, roomVisibility bool) ([]core.Room, error) {
	return s.repo.GetList(ctx, roomVisibility)
}

func (s *RoomsService) GetByID(ctx context.Context, roomID, userID primitive.ObjectID) (string, error) {
	r, err := s.repo.GetById(ctx, roomID)
	if err != nil {
		return "", err
	}

	invited := false

	for _, v := range r.Invites {
		if v == userID {
			invited = true
		}
	}

	if !r.Visibility && !invited {
		return "", core.ErrRoomNotFound
	}

	return s.room.GetJoinToken(roomID.String(), userID.String())
}
