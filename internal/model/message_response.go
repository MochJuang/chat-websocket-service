package model

type MessageResponse struct {
	MessageType string `json:"message_type"`
	Message     string `json:"message"`
}

const (
	MessageTypeChat         = "CHAT"
	MessageTypeNotification = "NOTIFICATION"
)
