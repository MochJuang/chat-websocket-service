package repository

type ChatRepository interface {
	GetConversation(conversationId int) ([]uint32, error)
	SendMessage(conversationId int, senderId int, content string) error
}
