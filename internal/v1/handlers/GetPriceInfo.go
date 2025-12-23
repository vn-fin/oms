package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/utils"
)

// GetPriceInfo returns bid1-3, ask1-3, mid, ceil, floor price for a symbol
// @Summary Get price info
// @Description Get bid1-3, ask1-3, mid, ceil, floor for a symbol from latest orderbook
// @Tags market
// @Param symbol path string true "Stock symbol"
// @Param price_level query string false "Price level: bid1, bid2, bid3, ask1, ask2, ask3, mid, ceil, floor (empty = all)"
// @Success 200 {object} controller.PriceInfo
// @Router /oms/v1/market/price/{symbol} [get]
func GetPriceInfo(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	if symbol == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "symbol is required",
		})
	}

	priceLevel := strings.ToLower(c.Query("price_level"))

	priceInfo, err := utils.GetPriceInfo(symbol)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":  "error fetching price data",
			"symbol": symbol,
			"detail": err.Error(),
		})
	}
	if priceInfo == nil {
		return c.Status(404).JSON(fiber.Map{
			"error":  "no price data for this symbol",
			"symbol": symbol,
		})
	}

	// If price_level specified, return only that level
	if priceLevel != "" {
		var price float64
		switch priceLevel {
		case "bid1":
			price = priceInfo.Bid1
		case "bid2":
			price = priceInfo.Bid2
		case "bid3":
			price = priceInfo.Bid3
		case "ask1":
			price = priceInfo.Ask1
		case "ask2":
			price = priceInfo.Ask2
		case "ask3":
			price = priceInfo.Ask3
		case "mid":
			price = priceInfo.Mid
		case "ceil":
			price = priceInfo.Ceil
		case "floor":
			price = priceInfo.Floor
		default:
			return c.Status(400).JSON(fiber.Map{
				"error":        "invalid price_level",
				"valid_levels": []string{"bid1", "bid2", "bid3", "ask1", "ask2", "ask3", "mid", "ceil", "floor"},
			})
		}
		return c.JSON(fiber.Map{
			"symbol":      symbol,
			"price_level": priceLevel,
			"price":       price,
		})
	}

	// Return all prices
	return c.JSON(fiber.Map{
		"symbol": symbol,
		"bid1":   priceInfo.Bid1,
		"bid2":   priceInfo.Bid2,
		"bid3":   priceInfo.Bid3,
		"ask1":   priceInfo.Ask1,
		"ask2":   priceInfo.Ask2,
		"ask3":   priceInfo.Ask3,
		"mid":    priceInfo.Mid,
		"ceil":   priceInfo.Ceil,
		"floor":  priceInfo.Floor,
	})
}
