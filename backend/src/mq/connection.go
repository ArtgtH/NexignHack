package mq

import (
	"backend/src/config"
	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)

func CreateConnection() *amqp.Connection {
	url := config.Config("RABBITMQ_URL")
	connection, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	return connection
}

func CreateChannel(connection *amqp.Connection) *amqp.Channel {
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	return channel
}

func CreateTaskQueue(ch *amqp.Channel) *amqp.Queue {
	q, err := ch.QueueDeclare(
		"task_q",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	return &q
}

func CreateResultQueue(ch *amqp.Channel) *amqp.Queue {
	q, err := ch.QueueDeclare(
		"result_q",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	return &q
}
