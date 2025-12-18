package conn

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog/log"
	"github.com/vn-fin/oms/internal/config"
)

func InitPriceStreamConsumer() error {
	var err error
	kafkaConfig := kafka.ConfigMap{
		"bootstrap.servers":          config.KafkaServers,
		"group.id":                   config.KafkaConsumerGroup,
		"auto.offset.reset":          "latest",
		"fetch.max.bytes":            52428800, // 50MB
		"max.partition.fetch.bytes":  10485760, // 10MB
		"queued.max.messages.kbytes": 65536,    // 64MB
		"enable.auto.commit":         false,
	}
	Consumer, err = kafka.NewConsumer(&kafkaConfig)
	if err != nil {
		return err
	}

	// Fetch metadata
	meta, err := Consumer.GetMetadata(nil, false, 5000)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to fetch metadata")
		return err
	}

	var partitions []kafka.TopicPartition
	topicMeta, exists := meta.Topics[config.KafkaMarketDataTopic]
	if !exists {
		return fmt.Errorf("metadata for topic %s not found", config.KafkaMarketDataTopic)
	}
	for _, p := range topicMeta.Partitions {
		log.Info().Msgf("Found partition %d for topic %s", p.ID, config.KafkaMarketDataTopic)
		partition := kafka.TopicPartition{
			Topic:     &config.KafkaMarketDataTopic,
			Partition: p.ID,
			Offset:    kafka.OffsetEnd, // Start at the end
		}
		partitions = append(partitions, partition)
	}

	log.Info().Msgf("Assigning %d partitions to consumer", len(partitions))
	// Directly assign partitions (bypasses group coordination)
	err = Consumer.Assign(partitions)
	if err != nil {
		log.Error().Err(err).Msg("Failed to assign partitions")
		return err
	}

	log.Info().Msgf("Kafka consumer manually assigned to topics: %s", config.KafkaMarketDataTopic)
	return nil

}
