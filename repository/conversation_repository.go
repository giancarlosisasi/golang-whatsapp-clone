package repository

import (
	"context"
	db "golang-whatsapp-clone/database/gen"
	customerrors "golang-whatsapp-clone/errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
)

type ConversationRepository interface {
	GetUserConversations(ctx context.Context, userID string) (*[]db.GetUserConversationsRow, error)
	GetLastMessageFromConversation(ctx context.Context, conversationID string) (*db.GetLastMessageRow, error)
	CreateConversation(ctx context.Context, conversationType string) (*db.Conversation, error)
	FindDirectConversation(ctx context.Context, user1ID string, user2ID string) (*db.Conversation, error)
}

type ConversationPostgresRepository struct {
	DBQueries *db.Queries
	logger    *zerolog.Logger
}

func NewConversationRepository(dbQueries *db.Queries, logger *zerolog.Logger) *ConversationPostgresRepository {
	return &ConversationPostgresRepository{
		DBQueries: dbQueries,
		logger:    logger,
	}
}

func (r *ConversationPostgresRepository) GetUserConversations(ctx context.Context, userID string) (*[]db.GetUserConversationsRow, error) {
	uId, err := fromStringToUUID(userID)
	if err != nil {
		return nil, customerrors.ErrInvalidUUIDValue

	}
	result, err := r.DBQueries.GetUserConversations(ctx, uId)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *ConversationPostgresRepository) GetLastMessageFromConversation(ctx context.Context, conversationID string) (*db.GetLastMessageRow, error) {
	uid, err := fromStringToUUID(conversationID)
	if err != nil {
		return nil, customerrors.ErrInvalidUUIDValue
	}

	ids := []pgtype.UUID{
		uid,
	}

	lastMessage, err := r.DBQueries.GetLastMessage(ctx, ids)
	if err != nil {
		return nil, customerrors.ErrResourceNotFound
	}

	if len(lastMessage) == 0 {
		return nil, nil
	}

	return &lastMessage[0], nil
}

func (r *ConversationPostgresRepository) CreateConversation(ctx context.Context, conversationType string) (*db.Conversation, error) {
	conversation, err := r.DBQueries.CreateConversation(ctx, conversationType)
	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (r *ConversationPostgresRepository) FindDirectConversation(ctx context.Context, user1ID string, user2ID string) (*db.Conversation, error) {
	u1Id, err := fromStringToUUID(user1ID)
	if err != nil {
		return nil, customerrors.ErrInvalidUUIDValue
	}

	u2Id, err := fromStringToUUID(user2ID)
	if err != nil {
		return nil, customerrors.ErrInvalidUUIDValue
	}

	existing, err := r.DBQueries.FindDirectConversation(ctx, db.FindDirectConversationParams{
		UserID:   u1Id,
		UserID_2: u2Id,
	})

	if err != nil {
		return nil, customerrors.ErrResourceNotFound
	}

	return &existing, nil
}
