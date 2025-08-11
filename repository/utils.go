package repository

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func fromStringToUUID(value string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return pgtype.UUID{}, err
	}

	pgUUID := pgtype.UUID{
		Bytes: parsed,
		Valid: true,
	}

	return pgUUID, nil
}
