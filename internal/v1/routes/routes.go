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

		// Market data
		api.Get("/market/price/:symbol", handlers.GetPriceInfo)

		// baskets
		api.Get("/baskets", authMiddleware, handlers.BasketList)
		api.Post("/baskets", authMiddleware, handlers.BasketCreate)

		// basket execute sessions (must be before /baskets/:id to avoid route conflict)
		api.Get("/baskets/execute-sessions", authMiddleware, handlers.BasketExecuteSessionsList)

		api.Get("/baskets/:id", authMiddleware, handlers.BasketDetail)
		api.Put("/baskets/:id", authMiddleware, handlers.BasketUpdate)
		api.Delete("/baskets/:id", authMiddleware, handlers.BasketDelete)
		api.Post("/baskets/:id/cancel", authMiddleware, handlers.BasketExecutionCancel)
		api.Post("/baskets/:basket_id/execute", authMiddleware, handlers.BasketExecute)
		api.Post("/baskets/:basket_id/executions/:execution_id/cancel", authMiddleware, handlers.BasketExecutionCancel)
		api.Post("/baskets/:basket_id/executions/:execution_id/match", authMiddleware, handlers.BasketFillAll)
		api.Get("/baskets/:basket_id/execute/:execution_id/orders", authMiddleware, handlers.UserOrderListBySession)
		api.Put("/baskets/:basket_id/execute/:execution_id/update-price", authMiddleware, handlers.UserOrderBatchUpdatePrice)

		// orders
		api.Put("/orders/:order_id/update", authMiddleware, handlers.UserOrderUpdatePrice)

		// credentials
		api.Get("/credentials", authMiddleware, handlers.CredentialListAll)
		api.Post("/credentials", authMiddleware, handlers.CredentialCreate)
		api.Put("/credentials/:credential_id", authMiddleware, handlers.CredentialUpdate)
		api.Delete("/credentials/:credential_id", authMiddleware, handlers.CredentialDelete)

		// credential groups
		api.Get("/credential-groups", authMiddleware, handlers.GetListGroup)
		api.Post("/credential-groups", authMiddleware, handlers.CredentialGroupCreate)
		api.Get("/credential-groups/:group_id/credentials", authMiddleware, handlers.CredentialListByGroup)
		api.Get("/credentials/:credential_id/groups", authMiddleware, handlers.GroupListByCredential)

	}
}
