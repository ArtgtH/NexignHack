package mq

import (
	"backend/src/config"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)

func RedisClient() *redis.Client {
	url := config.Config("REDIS_URL")
	opt, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)

	return rdb
}

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
	taskQueue := config.Config("RABBITMQ_TASK_QUEUE")
	q, err := ch.QueueDeclare(
		taskQueue,
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
