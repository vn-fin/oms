package stream

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/vn-fin/oms/internal/config"
	"github.com/vn-fin/oms/pkg/mem"
	"github.com/vn-fin/xpb/xpb/order"
)

// processMessageFromChan spawns multiple worker goroutines that process Kafka messages concurrently.
func processMessageFromChan(ctx context.Context, numWorkers int, onTick func(symbol string, tickBytes []byte)) {
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			log.Info().Int("worker", workerID).Msg("Starting Kafka message processor")

			for {
				select {
				case <-ctx.Done():
					log.Info().Int("worker", workerID).Msg("Shutting down Kafka message processor")
					return

				case msg, ok := <-chanKafkaMessage:
					if !ok {
						log.Warn().Int("worker", workerID).Msg("Channel closed, stopping processor")
						return
					}

					switch msg.MessageType {
					case config.MessageTypeOrderBook:
						var orderBookMessage order.OrderBookInfo
						if err := json.Unmarshal(msg.MessageBytes, &orderBookMessage); err != nil {
							log.Error().Err(err).Str("symbol", msg.Symbol).Msg("Failed to unmarshal OrderBook message")
							continue
						}
						// Build to parse timestamp
						orderBookMessage.Build()
						// Store latest order book
						mem.SetLatestOrderBook(msg.Symbol, mem.OrderBookInfo(orderBookMessage))

					case config.MessageTypeStockInfo:
						var stockInfoMessage mem.StockInfo
						if err := json.Unmarshal(msg.MessageBytes, &stockInfoMessage); err != nil {
							log.Error().Err(err).Str("symbol", msg.Symbol).Msg("Failed to unmarshal StockInfo message")
							continue
						}
						// Build to parse timestamp
						stockInfoMessage.Build()
						// Store latest stock info
						mem.SetLatestStockInfo(msg.Symbol, stockInfoMessage)

					default:
						log.Warn().Str("type", msg.MessageType).Msg("Unknown message type received")
					}
				}
			}
		}(i)
	}

	go func() {
		<-ctx.Done()
		wg.Wait()
		log.Info().Msg("All Kafka message processors stopped")
	}()
}
