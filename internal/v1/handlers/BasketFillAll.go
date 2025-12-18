package handlers

import "github.com/gofiber/fiber/v2"

// BasketFillAll
// @Summary Fill all basket execution
// @Description Fill all basket execution with match data
// @Tags Execution
// @Router /oms/v1/baskets/:basket_id/executions/{execution_id}/match [post]
func BasketFillAll(c *fiber.Ctx, credentialID string) error {

	return c.JSON(fiber.Map{
		"message": "Welcome to OMS API",
	})
}
