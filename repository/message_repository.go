package repository

import (
	db "golang-whatsapp-clone/database/gen"
)

type MessageRepository struct {
	DBQueries *db.Queries
}

func NewMessageRepository(dbQueries *db.Queries) *MessageRepository {
	return &MessageRepository{
		DBQueries: dbQueries,
	}
}

func (r *MessageRepository) CreateMessage() {

}
