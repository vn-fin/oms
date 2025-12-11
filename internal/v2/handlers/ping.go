package handlers

import "github.com/gofiber/fiber/v2"

// PingHandler is the handler for the ping route
// @Summary Ping route
// @Description Welcome route
// @Tags ping
// @Router /oms/v2/ping [get]
func PingHandler(c *fiber.Ctx) error {
	// Send a simple JSON response
	return c.JSON(fiber.Map{
		"message": "Welcome to OMS API",
	})
}
