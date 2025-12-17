package handlers

import "github.com/gofiber/fiber/v2"

// BasketFillSymbol
// @Summary Fill basket symbol
// @Description Fill basket symbol with fill data
// @Tags Execution
// @Router /oms/v1/baskets/:basket_id/symbols/:symbol/fill [post]
func BasketFillSymbol(c *fiber.Ctx) error {
	// Send a simple JSON response
	return c.JSON(fiber.Map{
		"message": "Welcome to OMS API",
	})
}
