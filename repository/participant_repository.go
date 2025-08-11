package repository

import (
	"context"
	db "golang-whatsapp-clone/database/gen"
	customerrors "golang-whatsapp-clone/errors"

	"github.com/rs/zerolog"
)

type ParticipantRepository interface {
	CreateParticipants(ctx context.Context, conversationID string, userIDs []string) error
}

type ParticipantPostgresRepository struct {
	dbQueries *db.Queries
	logger    *zerolog.Logger
}

func NewParticipantRepository(dbQueries *db.Queries, logger *zerolog.Logger) *ParticipantPostgresRepository {
	return &ParticipantPostgresRepository{
		dbQueries: dbQueries,
		logger:    logger,
	}
}

func (r *ParticipantPostgresRepository) CreateParticipants(ctx context.Context, conversationID string, userIDs []string) error {
	cId, err := fromStringToUUID(conversationID)
	if err != nil {
		r.logger.Error().Msgf("Repo:CreateParticipants: invalid conversation id: %v", err)
		return customerrors.ErrInvalidUUIDValue
	}

	var participants []db.CreateParticipantsParams
	for _, rawId := range userIDs {
		uId, err := fromStringToUUID(rawId)
		if err != nil {
			r.logger.Error().Msgf("Repo:CreateParticipants: invalid user 2 id, value is: %s and error is: %v", rawId, err)
			return customerrors.ErrInvalidUUIDValue
		}
		participants = append(participants, db.CreateParticipantsParams{
			ConversationID: cId,
			UserID:         uId,
			// TODO: add support for admin, owner, member
			Role: "member",
		})
	}

	_, err = r.dbQueries.CreateParticipants(ctx, participants)

	if err != nil {
		r.logger.Error().Msgf("Repo:CreateParticipants: error to bulk create participants, %v", err)
		return err
	}

	return nil
}
