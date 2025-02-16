package mocks

import (
	"backend/src/router/structs"
	"backend/src/service/messages"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRabbitRedisTaskService struct {
	mock.Mock
}

func (m *MockRabbitRedisTaskService) PublishTask(task *messages.CreatedFullTask) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockRabbitRedisTaskService) GetTaskResult(taskID uuid.UUID) (*structs.FileTaskResponse, error) {
	args := m.Called(taskID)
	return args.Get(0).(*structs.FileTaskResponse), args.Error(1)
}
