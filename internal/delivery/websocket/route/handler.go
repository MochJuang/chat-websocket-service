package route

import (
	"log"
	"strconv"
	"websocket-service/internal/config"
	wsdelivery "websocket-service/internal/delivery/websocket"
	"websocket-service/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutes(app *fiber.App, cfg config.Config) {

	manager := utils.NewWebSocketManager()

	go manager.Run()

	websocketController := wsdelivery.NewWebSocketController(manager)

	app.Use("/ws/:userId", websocket.New(func(c *websocket.Conn) {
		userIdStr := c.Params("userId")
		userId, err := strconv.Atoi(userIdStr)

		if err != nil {
			log.Printf("Invalid userId: %s", userIdStr)
			return
		}

		log.Println("New connection for user", userId)
		manager.WebSocketEndpoint(c, userId)
	}))

	app.Get("/ws/:userId", websocketController.Get)

	app.Get("/broadcast", websocketController.Broadcast)

}
