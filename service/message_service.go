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

func (s *MessageService) GetMessages(ctx context.Context, conversationID string, limit int32, offset int32) (*[]db.GetConversationMessagesRow, error) {
	result, err := s.MessageRepository.GetMessages(ctx, conversationID, limit, offset)

	if err != nil {
		return nil, err
	}

	return result, nil
}
