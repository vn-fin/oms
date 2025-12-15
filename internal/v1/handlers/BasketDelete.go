package handlers

import "github.com/gofiber/fiber/v2"

// BasketDelete
// @Summary
// @Description
// @Tags
// @Router /oms/v1/baskets [delete]
func BasketDelete(c *fiber.Ctx) error {
	// Send a simple JSON response
	return c.JSON(fiber.Map{
		"message": "Welcome to OMS API",
	})
}
