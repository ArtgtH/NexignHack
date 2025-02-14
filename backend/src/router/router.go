package router

import (
	"backend/src/router/handlers"
	"backend/src/service/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

// @title Songs Swagger API
// @version 1.0
// @description Swagger API for Nexign ML project.

// @BasePath /api/
// @host localhost:3000

func InitRouter(taskService services.TaskService) *fiber.App {
	router := fiber.New()
	router.Use(logger.New())
	router.Use(recover.New())
	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	binaryHandler := handlers.NewBinaryHandler(taskService)

	router.Get("/swagger/*", swagger.HandlerDefault)
	api := router.Group("/api")
	{
		api.Post("/binary/full/", binaryHandler.FullTask)
		api.Post("/binary/short/", binaryHandler.ShortTask)
	}

	return router
}
