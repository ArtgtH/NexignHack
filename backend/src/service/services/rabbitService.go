package services

import (
	"backend/src/router/structs"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type TaskService interface {
	PublishTask(task *structs.CreatedFullTask) error
	GetTaskResult(taskID uuid.UUID) (*structs.ResultTask, error)
}

type RabbitRedisTaskService struct {
	ch        *amqp.Channel
	taskQueue *amqp.Queue
	rdb       *redis.Client
	ctx       context.Context
}

func NewRabbitRedisTaskService(ch *amqp.Channel, taskQueue *amqp.Queue, rdb *redis.Client, ctx context.Context) *RabbitRedisTaskService {
	return &RabbitRedisTaskService{ch, taskQueue, rdb, ctx}
}

func (r *RabbitRedisTaskService) PublishTask(task *structs.CreatedFullTask) error {
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

	curr := time.Now()
	log.Println(curr)
	return nil
}

func (r *RabbitRedisTaskService) GetTaskResult(taskID uuid.UUID) (*structs.ResultTask, error) {
	key := taskID.String()

	for {
		val, err := r.rdb.Get(r.ctx, key).Result()
		if errors.Is(redis.Nil, err) {
			continue
		} else if err != nil {
			log.Fatal(err)
		} else {
			var res structs.ResultTask
			err = json.Unmarshal([]byte(val), &res)
			return &res, nil
		}
	}
}
