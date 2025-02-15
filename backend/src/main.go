package main

import (
	"backend/src/mq"
	"backend/src/router"
	"backend/src/service/services"
	"context"
	"log"
)

func main() {
	conn := mq.CreateConnection()
	ch := mq.CreateChannel(conn)
	taskQ := mq.CreateTaskQueue(ch)

	rdb := mq.RedisClient()
	ctx := context.Background()
	defer rdb.Close()
	defer ch.Close()
	defer conn.Close()

	rabbitService := services.NewRabbitTaskService(ch, taskQ, rdb, ctx)

	app := router.InitRouter(rabbitService)

	log.Fatal(app.Listen(":5050"))
}
