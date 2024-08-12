package model

import "github.com/gofiber/fiber/v2"

func Response(status string, message string, data interface{}) fiber.Map {
	return fiber.Map{
		"status":  status,
		"message": message,
		"data":    data,
	}
}
