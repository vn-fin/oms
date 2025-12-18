package utils

import (
	"fmt"

	"github.com/vn-fin/oms/internal/config"
)

// GenerateOrderATOKey

func GenerateOrderATOKey(symbol string) string {
	return fmt.Sprintf("%s.%s.ATO", config.RedisOrderListPrefixKey, symbol)
}

// GenerateOrderATCKey

func GenerateOrderATCKey(symbol string) string {
	return fmt.Sprintf("%s.%s.ATC", config.RedisOrderListPrefixKey, symbol)
}

// GenerateOrderListKey generates a Redis key for storing order lists based on symbol and price.
func GenerateOrderListKey(symbol string, side string, price float64) string {
	return fmt.Sprintf("%s.%s.%s.%.2f", config.RedisOrderListPrefixKey, symbol, side, price)
}

// GenerateLatestOrderKey generates a Redis key for storing an order by its ID.
func GenerateLatestOrderKey() string {
	return fmt.Sprintf("%s.%s", config.RedisOrderSetKey, "all")
}
