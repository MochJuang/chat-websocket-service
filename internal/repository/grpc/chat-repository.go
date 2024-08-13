package grpc

import (
	"context"
	"github.com/MochJuang/chat-grpc/service/chat"
	"google.golang.org/grpc"
	"log"
	e "websocket-service/internal/exception"
	"websocket-service/internal/repository"
)

type chatRepository struct {
	client           chat.ChatServiceClient
	conversationData map[int][]uint32
}

func NewChatRepository(client *grpc.ClientConn) repository.ChatRepository {
	return &chatRepository{
		client:           chat.NewChatServiceClient(client),
		conversationData: make(map[int][]uint32),
	}
}

func (r *chatRepository) GetConversation(conversationId int) ([]uint32, error) {

	log.Println("Conversations:", r.conversationData)
	if participants, ok := r.conversationData[conversationId]; ok {
		return participants, nil
	}

	req := &chat.ConversationRequest{
		ConversationId: uint32(conversationId),
	}

	res, err := r.client.GetConversationDetails(context.Background(), req)
	if err != nil {
		log.Println("Error getting conversation details:", err)
		return nil, e.NotFound("Conversation not found")
	}

	r.conversationData[conversationId] = res.ParticipantIds
	return res.ParticipantIds, nil
}

func (r *chatRepository) SendMessage(conversationId int, senderId int, content string) error {
	req := &chat.AddMessageRequest{
		ConversationId: uint32(conversationId),
		SenderId:       uint32(senderId),
		Content:        content,
	}

	_, err := r.client.AddMessageToConversation(context.Background(), req)
	if err != nil {
		return e.Internal(err)
	}

	return nil
}
