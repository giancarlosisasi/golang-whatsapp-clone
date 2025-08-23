package repository

import (
	"context"
	db "golang-whatsapp-clone/database/gen"
	customerrors "golang-whatsapp-clone/errors"

	"github.com/jackc/pgx/v5/pgtype"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, conversationID string, senderID string, content string, messageType string, replyToMessageID *string) (*db.Message, error)
	GetMessages(ctx context.Context, conversationID string, limit int32, offset int32) (*[]db.GetConversationMessagesRow, error)
}

type MessagePostgresRepository struct {
	DBQueries *db.Queries
}

func NewMessageRepository(dbQueries *db.Queries) *MessagePostgresRepository {
	return &MessagePostgresRepository{
		DBQueries: dbQueries,
	}
}

func (r *MessagePostgresRepository) CreateMessage(ctx context.Context, conversationID string, senderID string, content string, messageType string, replyToMessageID *string) (*db.Message, error) {
	cui, err := fromStringToUUID(conversationID)
	if err != nil {
		return nil, customerrors.ErrInvalidUUIDValue
	}

	sui, err := fromStringToUUID(senderID)
	if err != nil {
		return nil, customerrors.ErrInvalidUUIDValue
	}

	var rui pgtype.UUID

	if replyToMessageID != nil {
		rui, err = fromStringToUUID(*replyToMessageID)
		if err != nil {
			return nil, customerrors.ErrInvalidUUIDValue
		}
	}

	message, err := r.DBQueries.CreateMessage(ctx, db.CreateMessageParams{
		ConversationID:   cui,
		SenderID:         sui,
		Content:          content,
		MessageType:      messageType,
		ReplyToMessageID: rui,
		Status:           MESSAGE_STATUS_SENT,
	})

	if err != nil {
		return nil, err
	}

	return &message, nil

}

func (r *MessagePostgresRepository) GetMessages(ctx context.Context, conversationID string, limit int32, offset int32) (*[]db.GetConversationMessagesRow, error) {
	cui, err := fromStringToUUID(conversationID)
	if err != nil {
		return nil, customerrors.ErrInvalidUUIDValue
	}

	messages, err := r.DBQueries.GetConversationMessages(ctx, db.GetConversationMessagesParams{
		ConversationID: cui,
		Limit:          limit,
		Offset:         offset,
	})

	if err != nil {
		return nil, err
	}

	return &messages, nil
}
