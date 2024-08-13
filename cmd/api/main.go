package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"websocket-service/internal/config"
	"websocket-service/internal/delivery/websocket/route"
	"websocket-service/internal/utils"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// Initialize Fiber app
	app := fiber.New()

	app.Use(logger.New())

	rmq, err := utils.NewRabbitMQ(cfg.RabbitMQAddress)
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v", err)
	}
	cfg.RabbitMQUtils = rmq
	// Setup routes
	route.SetupRoutes(app, cfg)

	// Start server
	log.Fatal(app.Listen(cfg.ServerAddress))
}
