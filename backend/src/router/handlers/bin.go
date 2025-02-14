package handlers

import (
	"backend/src/router/structs"
	"backend/src/service/converter"
	"backend/src/service/services"
	"github.com/gofiber/fiber/v2"
)

type BinaryHandler struct {
	taskService services.TaskService
}

func NewBinaryHandler(taskService services.TaskService) *BinaryHandler {
	handler := BinaryHandler{
		taskService: taskService,
	}
	return &handler
}

// FullTask godoc
// @Summary Upload a file for a full task
// @Description Uploads a file and processes it for a full task
// @Accept multipart/form-data
// @Produce plain/text
// @Param file formData file true "File to upload"
// @Success 200 {string} string "File uploaded successfully"
// @Failure 400 {string} string "Error retrieving the file"
// @Failure 500 {string} string "Error opening or reading the file"
// @Router /binary/full [post]
func (h *BinaryHandler) FullTask(c *fiber.Ctx) error {
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

	task := structs.CreatedTask{
		Type:     "FullTask",
		Messages: data,
	}

	err = h.taskService.PublishTask(&task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error publishing the task")
	}

	return c.SendString("File uploaded successfully")
}

// ShortTask godoc
// @Summary Upload a file for a short task
// @Description Uploads a file and processes it for a short task
// @Accept multipart/form-data
// @Produce plain/text
// @Param file formData file true "File to upload"
// @Success 200 {string} string "File uploaded successfully"
// @Failure 400 {string} string "Error retrieving the file"
// @Failure 500 {string} string "Error opening or reading the file"
// @Router /binary/short [post]
func (h *BinaryHandler) ShortTask(c *fiber.Ctx) error {
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
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error working with file")
	}

	task := structs.CreatedTask{
		Type:     "ShortTask",
		Messages: data,
	}

	err = h.taskService.PublishTask(&task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error publishing the task")
	}

	return c.SendString("File uploaded successfully")
}
