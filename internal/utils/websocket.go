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
	broadcast  chan *BroadcastMessage
	register   chan *WebSocketConnInfo
	unregister chan *WebSocketConnInfo
	mu         sync.Mutex
}

type WebSocketConnInfo struct {
	UserId int
	Conn   *websocket.Conn
}

type BroadcastMessage struct {
	UserIds []uint32
	Message []byte
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients:    make(map[uint32][]*websocket.Conn),
		broadcast:  make(chan *BroadcastMessage),
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
		case broadcastMsg := <-manager.broadcast:
			manager.broadcastMessage(broadcastMsg)
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

func (manager *WebSocketManager) broadcastMessage(broadcastMsg *BroadcastMessage) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	log.Printf("Broadcasting message to %d users: %s", len(broadcastMsg.UserIds), string(broadcastMsg.Message))
	for _, userId := range broadcastMsg.UserIds {
		if conns, ok := manager.clients[userId]; ok {
			for _, conn := range conns {
				if err := conn.WriteMessage(websocket.TextMessage, broadcastMsg.Message); err != nil {
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

		manager.BroadcastMessageChat(userIds, request.Message)
		manager.BroadcastMessageNotification(userIds, request.Message)
	}
}

func (manager *WebSocketManager) BroadcastMessageChat(userIds []uint32, message string) {
	response := model.MessageResponse{
		MessageType: model.MessageTypeChat,
		Message:     message,
	}
	responseByte, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return
	}

	manager.broadcast <- &BroadcastMessage{
		UserIds: userIds,
		Message: responseByte,
	}
}

func (manager *WebSocketManager) BroadcastMessageNotification(userIds []uint32, message string) {
	response := model.MessageResponse{
		MessageType: model.MessageTypeNotification,
		Message:     message,
	}
	responseByte, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return
	}

	manager.broadcast <- &BroadcastMessage{
		UserIds: userIds,
		Message: responseByte,
	}

}

func (manager *WebSocketManager) BroadcastMessageToUsers(userIds []uint32, message []byte) {
	manager.broadcast <- &BroadcastMessage{
		UserIds: userIds,
		Message: message,
	}
}
