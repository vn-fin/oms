package stream

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vn-fin/oms/internal/config"
	"github.com/vn-fin/oms/pkg/controller"
	"github.com/vn-fin/oms/pkg/mem"
)

// Stats prints periodic Kafka statistics.
func Stats(ctx context.Context, printEverySeconds int) {
	ticker := time.NewTicker(time.Duration(printEverySeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Stats stopped: context cancelled")
			return
		case <-ticker.C:
			printStats()
		}
	}
}

func printStats() {
	currentReceived := totalReceivedMessages.Load()
	currentValid := totalValidMessages.Load()

	maxOrderBookTime := getMaxTimeFromMemory(config.MessageTypeOrderBook)
	pendingMessages := len(chanKafkaMessage)

	// Get price info for the latest symbol
	priceInfo := controller.GetPriceInfo(maxOrderBookTime.Symbol)

	logEvent := log.Info().
		Str("module", "stream-stats").
		Int64("total_messages", currentReceived).
		Str("latest_orderbook", maxOrderBookTime.ToString()).
		Int("pending_messages", pendingMessages).
		Int64("valid_messages", currentValid)

	if priceInfo != nil {
		logEvent = logEvent.
			Float64("bid1", priceInfo.Bid1).
			Float64("bid2", priceInfo.Bid2).
			Float64("bid3", priceInfo.Bid3).
			Float64("ask1", priceInfo.Ask1).
			Float64("ask2", priceInfo.Ask2).
			Float64("ask3", priceInfo.Ask3).
			Float64("mid", priceInfo.Mid)
	}

	logEvent.Msg("Kafka stream statistics")
}

func getMaxTimeFromMemory(dataType string) LatestMessage {
	switch dataType {
	case config.MessageTypeOrderBook:
		mem.Mutex.RLock()
		defer mem.Mutex.RUnlock()
		var maxTimeMessage LatestMessage
		for _, v := range mem.LatestOrderBookMap {
			currentMsg := LatestMessage{
				Time:   v.TimeF,
				Symbol: v.Symbol,
			}
			if currentMsg.Time > maxTimeMessage.Time {
				maxTimeMessage = currentMsg
			}
		}
		return maxTimeMessage
	default:
		log.Error().Str("dataType", dataType).Msg("Invalid data type for getMaxTimeFromMemory")
		return LatestMessage{}
	}
}
