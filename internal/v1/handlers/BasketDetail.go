package handlers

import "github.com/gofiber/fiber/v2"

// BasketDetail
// @Summary
// @Description
// @Tags
// @Router /oms/v1/baskets/:basket_id [get]
func BasketDetail(c *fiber.Ctx) error {
	// Send a simple JSON response
	return c.JSON(fiber.Map{
		"message": "Welcome to OMS API",
	})
}
