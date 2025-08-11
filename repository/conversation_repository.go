package repository

import (
	"context"
	db "golang-whatsapp-clone/database/gen"
	customerrors "golang-whatsapp-clone/errors"
)

type ConversationRepository interface {
	GetUserConversations(ctx context.Context, userID string) (*[]db.GetUserConversationsRow, error)
}

type ConversationPostgresRepository struct {
	DBQueries *db.Queries
}

func NewConversationRepository(dbQueries *db.Queries) *ConversationPostgresRepository {
	return &ConversationPostgresRepository{
		DBQueries: dbQueries,
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
