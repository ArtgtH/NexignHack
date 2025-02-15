package tests

import (
	"backend/src/router"
	"backend/src/router/structs"
	"backend/tests/mocks"
	"bytes"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mime/multipart"
	"net/http/httptest"
	"testing"
)

func TestBaseHandler_FullTask(t *testing.T) {
	app := fiber.New()
	mockService := new(mocks.MockRabbitRedisTaskService)
	router.InitRouter(mockService)

	t.Run("Success", func(t *testing.T) {
		expectedTask := &structs.CreatedFullTask{
			ID:   uuid.New(),
			Type: "FullTask",
		}
		expectedResult := &structs.ResultTask{
			ID: expectedTask.ID,
		}

		mockService.On("PublishTask", mock.Anything).Return(nil)
		mockService.On("GetTaskResult", expectedTask.ID).Return(expectedResult, nil)

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.xlsx")
		part.Write([]byte("test content"))
		writer.Close()

		req := httptest.NewRequest("POST", "/binary/full/", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Execute
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("Missing file", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/binary/full/", nil)
		req.Header.Set("Content-Type", "multipart/form-data")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}
