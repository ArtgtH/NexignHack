package main

import (
	_ "backend/src/docs"
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

	taskService := services.NewRabbitRedisTaskService(ch, taskQ, rdb, ctx)

	app := router.InitRouter(taskService)

	log.Fatal(app.Listen(":5050"))
}
