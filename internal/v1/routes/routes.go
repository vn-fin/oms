package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/config"
	"github.com/vn-fin/oms/internal/middlewares"
	"github.com/vn-fin/oms/internal/v1/handlers"
)

func SetupRoutes(app *fiber.App) {
	authMiddleware := middlewares.AuthMiddleware()

	api := app.Group(fmt.Sprintf("/%s/v1", config.ServiceName))
	{

		// Ping
		api.Get("/ping", handlers.PingHandler)
		api.Get("/ping-auth", authMiddleware, handlers.PingHandler)

		// baskets
		api.Get("/baskets", authMiddleware, handlers.BasketList)
		api.Post("/baskets", authMiddleware, handlers.BasketCreate)
		api.Get("/baskets/:id", authMiddleware, handlers.BasketDetail)
		api.Delete("/baskets/:id", authMiddleware, handlers.BasketDelete)
		api.Post("/baskets/:id/cancel", authMiddleware, handlers.BasketCancel)

		// basket execute sessions
		api.Post("/baskets/:basket_id/execute", authMiddleware, handlers.BasketExecute)
		api.Get("/baskets/:basket_id/execute-sessions", authMiddleware, handlers.BasketExecuteSessionsList)
		api.Post("/baskets/:basket_id/executions/:execution_id/cancel", authMiddleware, handlers.BasketExecutionCancel)

	}
}
