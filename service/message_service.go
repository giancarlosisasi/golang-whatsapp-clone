package service

import "golang-whatsapp-clone/repository"

type MessageService struct {
	MessageRepository *repository.MessageRepository
}

func NewMessageService(messageRepository *repository.MessageRepository) *MessageService {
	return &MessageService{
		MessageRepository: messageRepository,
	}
}
