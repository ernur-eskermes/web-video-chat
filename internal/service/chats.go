package service

import (
	"context"
	"github.com/ernur-eskermes/web-video-chat/internal/domain"
	"github.com/ernur-eskermes/web-video-chat/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/olahol/melody.v1"
	"time"
)

type ChatsService struct {
	repo      repository.Chats
	websocket *melody.Melody
}

func NewChatsService(repo repository.Chats, websocket *melody.Melody) *ChatsService {
	return &ChatsService{
		repo:      repo,
		websocket: websocket,
	}
}

func (c *ChatsService) GetMessages(ctx context.Context, id primitive.ObjectID) ([]domain.Message, error) {
	return c.repo.GetChatMessages(ctx, id)
}

func (c *ChatsService) CreateMessage(ctx context.Context, input CreateMessageInput) error {
	return c.repo.CreateMessage(ctx, domain.Message{
		Sender:    input.UserId,
		ChatId:    input.ChatId,
		Text:      input.Message,
		CreatedAt: time.Now(),
	})
}
