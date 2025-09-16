package broker

import (
	"context"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	rabbit *RabbitMQ
	queue  string
}

func NewProducer(r *RabbitMQ, queue string) *Producer {
	_, err := r.Channel.QueueDeclare(
		queue,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		panic(err)
	}

	return &Producer{rabbit: r, queue: queue}
}

func (p *Producer) Publish(message []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return p.rabbit.Channel.PublishWithContext(ctx,
		"",      // exchange
		p.queue, // routing key
		false,   // mandatory
		false,   // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        message,
		})
}

// func (p *Producer) ExtendChannel(subqueue string) *Producer {
// 		_, err := p.rabbit.Channel.QueueDeclare(
// 		fmt.Sprintf("%s-%d", p.queue, subqueue) ,
// 		true,  // durable
// 		false, // autoDelete
// 		false, // exclusive
// 		false, // noWait
// 		nil,
// 	)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return &Producer{rabbit: p.rabbit, queue: queue}
// }
