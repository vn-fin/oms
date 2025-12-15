package handlers

import "github.com/gofiber/fiber/v2"

// BasketCreate
// @Summary
// @Description
// @Tags
// @Router /oms/v1/baskets [post]
func BasketCreate(c *fiber.Ctx) error {
	// Send a simple JSON response
	return c.JSON(fiber.Map{
		"message": "Welcome to OMS API",
	})
}
