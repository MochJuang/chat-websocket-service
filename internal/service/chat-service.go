package service

import (
	e "websocket-service/internal/exception"
	"websocket-service/internal/model"
	"websocket-service/internal/repository"
	"websocket-service/internal/utils"
)

type ChatService interface {
	SendMessage(conversationId int, senderId int, content string) error
	GetConversation(conversationId int) ([]uint32, error)
	ProcessMessage(userId int, request model.MessageRequest) ([]uint32, error)
}

type chatService struct {
	repo repository.ChatRepository
}

func NewChatService(repo repository.ChatRepository) ChatService {
	return &chatService{repo: repo}
}

func (s *chatService) ProcessMessage(userId int, request model.MessageRequest) ([]uint32, error) {
	err := utils.Validate(request)
	if err != nil {
		return []uint32{}, e.Validation(err)
	}

	err = s.SendMessage(request.ConversationId, userId, request.Message)
	if err != nil {
		return []uint32{}, err
	}

	userIds, err := s.GetConversation(request.ConversationId)
	if err != nil {
		return []uint32{}, err
	}

	return userIds, nil

}

func (s *chatService) SendMessage(conversationId int, senderId int, content string) error {
	err := s.repo.SendMessage(conversationId, senderId, content)
	if err != nil {
		return err
	}
	return nil
}

func (s *chatService) GetConversation(conversationId int) ([]uint32, error) {
	messages, err := s.repo.GetConversation(conversationId)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
