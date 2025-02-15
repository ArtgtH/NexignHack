package mocks

import (
	"backend/src/router/structs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRabbitRedisTaskService struct {
	mock.Mock
}

func (m *MockRabbitRedisTaskService) PublishTask(task *structs.CreatedFullTask) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockRabbitRedisTaskService) GetTaskResult(taskID uuid.UUID) (*structs.ResultTask, error) {
	args := m.Called(taskID)
	return args.Get(0).(*structs.ResultTask), args.Error(1)
}
