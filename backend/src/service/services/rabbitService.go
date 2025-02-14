package services

import (
	"backend/src/router/structs"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type TaskService interface {
	PublishTask(task *structs.CreatedTask) error
	GetTaskResult(taskID uint) (string, error)
}

type RabbitTaskService struct {
	ch          *amqp.Channel
	taskQueue   *amqp.Queue
	resultQueue *amqp.Queue
}

func NewRabbitTaskService(ch *amqp.Channel, taskQueue *amqp.Queue, resultQueue *amqp.Queue) *RabbitTaskService {
	return &RabbitTaskService{ch, taskQueue, resultQueue}
}

func (r *RabbitTaskService) PublishTask(task *structs.CreatedTask) error {
	body, err := json.Marshal(task)
	if err != nil {
		return err
	}

	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}

	err = r.ch.Publish(
		"",
		r.taskQueue.Name,
		false,
		false,
		message,
	)
	if err != nil {
		return err
	}

	log.Printf("Сообщение отправлено в очередь: %s", body)
	return nil
}

func (r *RabbitTaskService) GetTaskResult(taskID uint) (string, error) {
	return "nil", nil
}
