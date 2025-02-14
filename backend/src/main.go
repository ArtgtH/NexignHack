package main

import (
	"backend/src/mq"
	"backend/src/router"
	"backend/src/service/services"
	"log"
)

func main() {
	conn := mq.CreateConnection()
	ch := mq.CreateChannel(conn)
	taskQ := mq.CreateTaskQueue(ch)
	resQ := mq.CreateResultQueue(ch)

	rabbitService := services.NewRabbitTaskService(ch, taskQ, resQ)

	app := router.InitRouter(rabbitService)

	log.Fatal(app.Listen(":3000"))
}
