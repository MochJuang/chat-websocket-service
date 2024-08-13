package rabbitmq

import (
	"encoding/json"
	"log"
	"websocket-service/internal/entity"
	"websocket-service/internal/service"
	"websocket-service/internal/utils"
)

type RabbitMQConsumer struct {
	chatService   *service.ChatService
	manager       *utils.WebSocketManager
	rabbitMQUtils *utils.RabbitMQ
}

func NewRabbitMQConsumer(manager *utils.WebSocketManager, rabbitMQUtils *utils.RabbitMQ) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		manager:       manager,
		rabbitMQUtils: rabbitMQUtils,
	}
}

func (r *RabbitMQConsumer) Run() {
	go r.StartConsumeNotification()
	go r.StartConsumeBroadcast()
}

func (r *RabbitMQConsumer) StartConsumeNotification() {
	err := r.rabbitMQUtils.DeclareQueue(utils.QUEUE_NOTIFICATION)
	if err != nil {
		log.Fatalf("Failed to declare queue2: %v", err)
	}

	err = r.rabbitMQUtils.ConsumeMessages(utils.QUEUE_NOTIFICATION, "consumer1", func(body string) {
		var notification entity.Notification
		err = json.Unmarshal([]byte(body), &notification)
		if err != nil {
			log.Fatalf("Failed to unmarshal notification: %v", err)
		}

		r.manager.JobMessageNotification([]uint32{uint32(notification.UserID)}, notification.Message)

	})
	if err != nil {
		log.Fatalf("Failed to consume messages from queue1: %v", err)
	}

}

func (r *RabbitMQConsumer) StartConsumeBroadcast() {
	err := r.rabbitMQUtils.DeclareQueue(utils.QUEUE_BROADCAST)
	if err != nil {
		log.Fatalf("Failed to declare queue2: %v", err)
	}

	err = r.rabbitMQUtils.ConsumeMessages(utils.QUEUE_BROADCAST, "consumer2", func(body string) {
		var notification entity.Notification
		err = json.Unmarshal([]byte(body), &notification)
		if err != nil {
			log.Fatalf("Failed to unmarshal notification: %v", err)
		}

		r.manager.BroadcastNotification(notification.Message)

	})
	if err != nil {
		log.Fatalf("Failed to consume messages from queue1: %v", err)
	}

}
