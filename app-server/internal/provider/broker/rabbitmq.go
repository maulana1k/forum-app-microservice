package broker

import (
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp091.Connection
	Channel *amqp091.Channel
	URL     string
}

func NewRabbitMQ(url string) *RabbitMQ {
	r := &RabbitMQ{URL: url}
	r.connect()
	return r
}

func (r *RabbitMQ) connect() {
	var err error
	for i := range 5 {
		r.Conn, err = amqp091.Dial(r.URL)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Cannot connect to RabbitMQ: %v", err)
	}

	r.Channel, err = r.Conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
}

func (r *RabbitMQ) Close() {
	if r.Channel != nil {
		_ = r.Channel.Close()
	}
	if r.Conn != nil {
		_ = r.Conn.Close()
	}
}

func (r *RabbitMQ) DeclareQueue(queueName string) (amqp091.Queue, error) {
	return r.Channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
}
