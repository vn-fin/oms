package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog/log"
	"github.com/vn-fin/oms/internal/config"
)

var Consumer *kafka.Consumer

// InitConsumer initializes the Kafka consumer2
func InitConsumer() error {
	var err error
	Consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.KafkaServers,
		"group.id":          config.KafkaConsumerGroup,
		"auto.offset.reset": "latest",
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to create Kafka consumer")
		return err
	}

	// Subscribe to topic
	err = Consumer.SubscribeTopics([]string{config.KafkaMarketDataTopic}, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to subscribe to Kafka topic")
		return err
	}

	log.Info().
		Str("servers", config.KafkaServers).
		Str("group", config.KafkaConsumerGroup).
		Str("topic", config.KafkaMarketDataTopic).
		Msg("Kafka consumer initialized successfully")

	return nil
}

// CloseConsumer closes the Kafka consumer
func CloseConsumer() {
	if Consumer != nil {
		Consumer.Close()
		log.Info().Msg("Kafka consumer closed")
	}
}
