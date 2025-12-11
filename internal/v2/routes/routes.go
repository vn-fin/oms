package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/config"
	"github.com/vn-fin/oms/internal/middlewares"
	"github.com/vn-fin/oms/internal/v2/handlers"
)

func SetupRoutes(app *fiber.App) {
	authMiddleware := middlewares.AuthMiddleware()

	api := app.Group(fmt.Sprintf("/%s/v2", config.ServiceName))
	{

		// Ping
		api.Get("/ping", handlers.PingHandler)
		api.Get("/ping-auth", authMiddleware, handlers.PingHandler)

	}
}
