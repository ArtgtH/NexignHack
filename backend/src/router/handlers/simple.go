package handlers

import (
	"backend/src/router/structs"
	"backend/src/service/converter"
	"backend/src/service/messages"
	_ "backend/src/service/messages"
	"backend/src/service/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"time"
)

type BaseHandler struct {
	taskService services.TaskService
}

func NewBaseHandler(taskService services.TaskService) *BaseHandler {
	handler := BaseHandler{
		taskService: taskService,
	}
	return &handler
}

// FullTask godoc
// @Summary Upload a file for a full task
// @Description Uploads a file and processes it for a task
// @Accept multipart/form-data
// @Produce plain/text
// @Param file formData file true "File to upload"
// @Success 201 {object} structs.FileTaskResponse
// @Failure 400 {string} string "Error retrieving the file"
// @Failure 500 {string} string "Error opening or reading the file"
// @Router /binary/full/ [post]
func (h *BaseHandler) FullTask(c *fiber.Ctx) error {

	curr := time.Now()
	log.Println("Приняли запрос:", curr)

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error retrieving the file")
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error opening the file")
	}
	defer src.Close()

	data, err := converter.ConvertFromXLSX(src)

	task := messages.CreatedFullTask{
		ID:       uuid.New(),
		Type:     "FullTask",
		Messages: data,
	}

	err = h.taskService.PublishTask(&task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error publishing the task")
	}

	curr = time.Now()
	log.Println("Отправили таску:", curr)

	res, err := h.taskService.GetTaskResult(task.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving the task")
	}
	curr = time.Now()
	log.Println("Взяли результат:", curr)

	return c.Status(fiber.StatusCreated).JSON(res)
}

// ShortTask godoc
// @Summary Upload a text for a short task
// @Description Uploads a text and processes it for a task
// @Accept json
// @Produce json
// @Param text body structs.TextTaskRequest true "Text for a task"
// @Success 201 {object} structs.TextTaskResponse
// @Failure 400 {string} string "Error retrieving the file"
// @Failure 500 {string} string "Error opening or reading the file"
// @Router /binary/short [post]
func (h *BaseHandler) ShortTask(c *fiber.Ctx) error {
	var text structs.TextTaskRequest
	if err := c.BodyParser(&text); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error retrieving the task")
	}

	msg := []messages.Message{
		{
			UserID:      "1",
			SubmitDate:  "1",
			MessageText: text.Text,
		},
	}

	task := messages.CreatedFullTask{
		ID:       uuid.New(),
		Type:     "ShortTask",
		Messages: msg,
	}

	err := h.taskService.PublishTask(&task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error publishing the task")
	}

	res, err := h.taskService.GetTaskResult(task.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving the task")
	}

	response := structs.TextTaskResponse{
		Text:   text.Text,
		Result: res.Messages[0].Result,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
