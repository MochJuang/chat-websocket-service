package route

import (
	"google.golang.org/grpc"
	"log"
	"websocket-service/internal/config"
	wsdelivery "websocket-service/internal/delivery/websocket"
	"websocket-service/internal/service"
	"websocket-service/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	grpcrepository "websocket-service/internal/repository/grpc"
)

func SetupRoutes(app *fiber.App, cfg config.Config) {

	conn, err := grpc.Dial(cfg.GrpcServer, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	//defer conn.Close()

	cfg.GrpcClient = conn

	manager := utils.NewWebSocketManager()

	go manager.Run()

	chatRepository := grpcrepository.NewChatRepository(cfg.GrpcClient)
	chatService := service.NewChatService(chatRepository)

	websocketController := wsdelivery.NewWebSocketController(manager, chatService)

	app.Use("/ws/:userId", websocket.New(websocketController.Connect))

	app.Get("/ws/:userId", websocketController.Get)

	app.Get("/broadcast", websocketController.Broadcast)

}
