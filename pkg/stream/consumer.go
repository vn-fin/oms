package stream

import (
	"context"
	"encoding/json"
	"sync/atomic"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog/log"
	"github.com/vn-fin/oms/internal/config"
	kafkaConn "github.com/vn-fin/oms/pkg/kafka"
)

var dropCounter atomic.Int64

// ProcessKafkaPriceStream streams market data from Kafka and processes it.
// onTick is an optional callback that gets called when a tick message is processed.
func ProcessKafkaPriceStream(ctx context.Context, onTick func(symbol string, tickBytes []byte)) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("Ctx Done! StreamMarketDataFromKafka received shutdown signal")
				return
			default:
				ev := kafkaConn.Consumer.Poll(100)
				if ev == nil {
					continue
				}

				switch e := ev.(type) {
				case *kafka.Message:
					totalReceivedMessages.Add(1)
					var msg KafkaMessage
					if err := json.Unmarshal(e.Value, &msg); err != nil {
						log.Error().Err(err).Msg("Failed to unmarshal Kafka message")
						continue
					}

					if msg.Source != config.KafkaMessageSource || !config.KafkaValidDataTypes[msg.MessageType] {
						continue
					}

					msg.Update(e.Value)
					select {
					case chanKafkaMessage <- msg:
					default:
						if dropCounter.Add(1)%1000 == 0 {
							log.Warn().Str("channel", "chanKafkaMessage").Msg("Channel full — dropping message")
						}
					}

					totalValidMessages.Add(1)
				case kafka.Error:
					log.Error().Msgf("Kafka error: %v", e)
				}
			}
		}
	}()

	go processMessageFromChan(ctx, 10, onTick)
}
