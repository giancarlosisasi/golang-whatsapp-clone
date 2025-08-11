package service

import (
	"context"
	"errors"
	db "golang-whatsapp-clone/database/gen"
	"golang-whatsapp-clone/repository"
)

type ConversationService struct {
	conversationRepository repository.ConversationRepository
}

func NewConversationService(conversationRepository repository.ConversationRepository) *ConversationService {
	return &ConversationService{
		conversationRepository: conversationRepository,
	}
}

func (s *ConversationService) GetUserConversations(ctx context.Context, userID string) (*[]db.GetUserConversationsRow, error) {
	result, err := s.conversationRepository.GetUserConversations(ctx, userID)
	if err != nil {
		return nil, errors.New("invalid user id value")
	}

	return result, nil
}
