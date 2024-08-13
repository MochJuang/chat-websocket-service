package utils

import (
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queues  map[string]amqp.Queue
}

const (
	QUEUE_NOTIFICATION = "notification"
	QUEUE_BROADCAST    = "broadcast"
)

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
		queues:  make(map[string]amqp.Queue),
	}, nil
}

func (r *RabbitMQ) DeclareQueue(queueName string) error {
	q, err := r.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	r.queues[queueName] = q
	return nil
}

func (r *RabbitMQ) PublishMessage(queueName, body string) error {
	q, exists := r.queues[queueName]
	if !exists {
		return amqp.ErrClosed // or a custom error indicating queue doesn't exist
	}

	err := r.channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		return err
	}
	log.Printf(" [x] Sent to %s: %s", queueName, body)
	return nil
}

func (r *RabbitMQ) ConsumeMessages(queueName, consumerName string, handler func(string)) error {
	q, exists := r.queues[queueName]
	if !exists {
		return amqp.ErrClosed // or a custom error indicating queue doesn't exist
	}

	msgs, err := r.channel.Consume(
		q.Name,       // queue
		consumerName, // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message from %s: %s", queueName, d.Body)
			handler(string(d.Body))
		}
	}()
	return nil
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}
