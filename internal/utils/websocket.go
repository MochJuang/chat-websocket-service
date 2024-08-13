package utils

import (
	"encoding/json"
	"log"
	"sync"
	"websocket-service/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type WebSocketManager struct {
	clients    map[uint32][]*websocket.Conn
	job        chan *jobMessage
	register   chan *WebSocketConnInfo
	unregister chan *WebSocketConnInfo
	mu         sync.Mutex
}

type WebSocketConnInfo struct {
	UserId int
	Conn   *websocket.Conn
}

type jobMessage struct {
	UserIds []uint32
	Message []byte
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients:    make(map[uint32][]*websocket.Conn),
		job:        make(chan *jobMessage),
		register:   make(chan *WebSocketConnInfo),
		unregister: make(chan *WebSocketConnInfo),
	}
}

func (manager *WebSocketManager) Run() {
	for {
		select {
		case connInfo := <-manager.register:
			manager.addClient(connInfo.UserId, connInfo.Conn)
		case connInfo := <-manager.unregister:
			manager.removeClient(uint32(connInfo.UserId), connInfo.Conn)
		case jobMsg := <-manager.job:
			manager.jobMessage(jobMsg)
		}
	}
}

func (manager *WebSocketManager) addClient(userId int, conn *websocket.Conn) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	manager.clients[uint32(userId)] = append(manager.clients[uint32(userId)], conn)
	log.Printf("New connection for user %d: %s", userId, conn.RemoteAddr().String())
}

func (manager *WebSocketManager) removeClient(userId uint32, conn *websocket.Conn) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if conns, ok := manager.clients[userId]; ok {
		for i, c := range conns {
			if c == conn {
				manager.clients[userId] = append(conns[:i], conns[i+1:]...)
				conn.Close()
				log.Printf("Connection closed for user %d", userId)
				break
			}
		}

		if len(manager.clients[userId]) == 0 {
			delete(manager.clients, userId)
		}
	}
}

func (manager *WebSocketManager) jobMessage(jobMsg *jobMessage) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	log.Printf("jobing message to %d users: %s", len(jobMsg.UserIds), string(jobMsg.Message))
	for _, userId := range jobMsg.UserIds {
		if conns, ok := manager.clients[userId]; ok {
			for _, conn := range conns {
				if err := conn.WriteMessage(websocket.TextMessage, jobMsg.Message); err != nil {
					log.Printf("Failed to send message to user %d at %s: %v", userId, conn.RemoteAddr().String(), err)
					manager.removeClient(userId, conn)
				}
			}
		} else {
			log.Printf("User %d is not connected", userId)
		}
	}
}

// public functions areas

func (manager *WebSocketManager) HandleWebSocket(c *fiber.Ctx) error {

	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (manager *WebSocketManager) WebSocketEndpoint(c *websocket.Conn, userId int, callback func(userId int, request model.MessageRequest) ([]uint32, error)) {
	connInfo := &WebSocketConnInfo{
		UserId: userId,
		Conn:   c,
	}

	manager.register <- connInfo

	defer func() {
		manager.unregister <- connInfo
	}()

	for {
		_, message, err := c.ReadMessage()

		var request model.MessageRequest
		err = json.Unmarshal(message, &request)
		if err != nil {
			log.Printf("request format invalid %d: %v", userId, err)
			continue
		}

		userIds, err := callback(userId, request)
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			continue
		}

		log.Printf("Received message from user %d: %s", userId, message)

		manager.JobMessageChat(userIds, request.Message)
		manager.JobMessageNotification(userIds, request.Message)
	}
}

func (manager *WebSocketManager) BroadcastNotification(message string) {
	var userIds []uint32
	for userId := range manager.clients {
		userIds = append(userIds, userId)
	}

	manager.JobMessageNotification(userIds, message)
}

func (manager *WebSocketManager) JobMessageChat(userIds []uint32, message string) {
	response := model.MessageResponse{
		MessageType: model.MessageTypeChat,
		Message:     message,
	}
	responseByte, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return
	}

	manager.job <- &jobMessage{
		UserIds: userIds,
		Message: responseByte,
	}
}

func (manager *WebSocketManager) JobMessageNotification(userIds []uint32, message string) {
	response := model.MessageResponse{
		MessageType: model.MessageTypeNotification,
		Message:     message,
	}
	responseByte, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return
	}

	manager.job <- &jobMessage{
		UserIds: userIds,
		Message: responseByte,
	}

}

func (manager *WebSocketManager) jobMessageToUsers(userIds []uint32, message []byte) {
	manager.job <- &jobMessage{
		UserIds: userIds,
		Message: message,
	}
}
