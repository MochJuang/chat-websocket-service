package model

type MessageRequest struct {
	Message        string `json:"message"`
	ConversationId int    `json:"conversation_id"`
}

var DummyConversation = map[int][]int{
	1: {2, 3},
	2: {1, 3},
	3: {1, 2},
	4: {1, 2, 3},
}
