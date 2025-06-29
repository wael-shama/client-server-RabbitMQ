package queue

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

const (
	QueueName = "command_queue"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		QueueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQ{conn, ch}, nil
}

func (rmq *RabbitMQ) Publish(message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return rmq.ch.Publish(
		"",        // exchange
		QueueName, // key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (rmq *RabbitMQ) Consume() (<-chan amqp.Delivery, error) {
	return rmq.ch.Consume(
		QueueName,
		"",    // consumer
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
}

func (rmq *RabbitMQ) Close() {
	rmq.ch.Close()
	rmq.conn.Close()
}
