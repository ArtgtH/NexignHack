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
// @termsOfService http://swagger.io/terms/

// @BasePath /ai
// @host localhost:5050

func InitRouter(taskService services.TaskService) *fiber.App {
	router := fiber.New()
	router.Use(logger.New())
	router.Use(recover.New())
	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	baseHandler := handlers.NewBaseHandler(taskService)

	router.Get("/swagger/*", swagger.HandlerDefault)
	api := router.Group("/ai")
	{
		api.Post("/binary/full/", baseHandler.FullTask)
		api.Post("/binary/short/", baseHandler.ShortTask)
	}

	return router
}
