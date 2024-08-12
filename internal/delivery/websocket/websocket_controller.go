package websocket

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
	"strconv"
	e "websocket-service/internal/exception"
	"websocket-service/internal/model"
	"websocket-service/internal/utils"
)

type WebSocketController struct {
	manager *utils.WebSocketManager
}

func NewWebSocketController(manager *utils.WebSocketManager) *WebSocketController {
	return &WebSocketController{}
}

func (controller *WebSocketController) Connect(c *websocket.Conn) {
	userIdStr := c.Params("userId")
	userId, err := strconv.Atoi(userIdStr)

	if err != nil {
		log.Printf("Invalid userId: %s", userIdStr)
		return
	}

	log.Println("New connection for user", userId)
	controller.manager.WebSocketEndpoint(c, userId)
}

func (controller *WebSocketController) Get(c *fiber.Ctx) error {
	userIdStr := c.Params("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		log.Printf("Invalid userId: %s", userIdStr)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid userId")
	}

	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("userId", userId)
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}

func (controller *WebSocketController) Broadcast(c *fiber.Ctx) error {
	var request model.MessageRequest
	if err := c.BodyParser(&request); err != nil {
		return e.Validation(err)
	}

	userIds := model.DummyConversation[request.ConversationId]
	message := []byte(request.Message)

	controller.manager.BroadcastMessageToUsers(userIds, message)
	return c.SendString("Broadcast sent")
}
