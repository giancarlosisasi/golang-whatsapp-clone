package service

import (
	"context"
	"errors"
	db "golang-whatsapp-clone/database/gen"
	"golang-whatsapp-clone/repository"
)

type MessageService struct {
	MessageRepository repository.MessageRepository
}

func NewMessageService(messageRepository repository.MessageRepository) *MessageService {
	return &MessageService{
		MessageRepository: messageRepository,
	}
}

func (s *MessageService) CreateMessage(ctx context.Context, conversationID string, senderID string, content string, messageType string, replyToMessageID *string) (*db.Message, error) {
	message, err := s.MessageRepository.CreateMessage(ctx, conversationID, senderID, content, messageType, replyToMessageID)

	if err != nil {
		return nil, errors.New("here was an error when creating the message")
	}

	return message, nil
}
