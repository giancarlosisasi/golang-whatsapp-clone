package service

import (
	"context"
	"errors"
	db "golang-whatsapp-clone/database/gen"
	"golang-whatsapp-clone/graph/model"
	"golang-whatsapp-clone/repository"
)

type ConversationService struct {
	conversationRepository repository.ConversationRepository
	participantRepository  repository.ParticipantRepository
}

func NewConversationService(
	conversationRepository repository.ConversationRepository,
	participantRepository repository.ParticipantRepository,
) *ConversationService {
	return &ConversationService{
		conversationRepository: conversationRepository,
		participantRepository:  participantRepository,
	}
}

func (s *ConversationService) GetUserConversations(ctx context.Context, userID string) (*[]db.GetUserConversationsRow, error) {
	result, err := s.conversationRepository.GetUserConversations(ctx, userID)
	if err != nil {
		return nil, errors.New("invalid user id value")
	}

	return result, nil
}

func (s *ConversationService) GetLastMessageFromConversation(ctx context.Context, conversationID string) (*db.GetLastMessageRow, error) {
	result, err := s.conversationRepository.GetLastMessageFromConversation(ctx, conversationID)
	if err != nil {
		return nil, errors.New("error to get the last message")
	}

	return result, nil
}

func (s *ConversationService) GetOrCreateDirectConversation(ctx context.Context, user1ID string, user2ID string) (*db.Conversation, error) {
	// first try to find existing conversation
	existing, err := s.conversationRepository.FindDirectConversation(ctx, user1ID, user2ID)
	if err == nil {
		return existing, nil
	}

	// create new conversation
	conversation, err := s.conversationRepository.CreateConversation(ctx, model.ConversationTypeEnumDirect.String())
	if err != nil {
		return nil, err
	}

	// add participants to the new created conversation
	err = s.participantRepository.CreateParticipants(ctx, conversation.ID.String(), []string{user1ID, user2ID})

	if err != nil {
		return nil, err
	}

	return conversation, nil
}
