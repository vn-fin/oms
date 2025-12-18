package conn

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog/log"
	"github.com/vn-fin/oms/internal/config"
)

var Producer *kafka.Producer

// InitKafkaProducer initializes Kafka producer connection
func InitKafkaProducer() error {
	var err error
	kafkaConfig := kafka.ConfigMap{
		"bootstrap.servers":  config.KafkaServers,
		"acks":               "all",
		"compression.type":   "lz4",
		"retries":            3,
		"linger.ms":          5,
		"batch.num.messages": 1000,
	}

	Producer, err = kafka.NewProducer(&kafkaConfig)
	if err != nil {
		log.Error().Err(err).Msg("failed to create kafka producer")
		return err
	}

	log.Info().Msgf("Kafka producer initialized at %s", config.KafkaServers)
	return nil
}

// ProduceKafka sends a message to Kafka topic
func ProduceKafka(topic string, message []byte) error {
	if Producer == nil {
		return kafka.NewError(kafka.ErrState, "producer not initialized", false)
	}

	deliveryChan := make(chan kafka.Event, 1)
	err := Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, deliveryChan)
	if err != nil {
		return err
	}

	// Wait for delivery report
	e := <-deliveryChan
	m := e.(*kafka.Message)
	close(deliveryChan)

	if m.TopicPartition.Error != nil {
		log.Error().Err(m.TopicPartition.Error).Str("topic", topic).Msg("failed to deliver kafka message")
		return m.TopicPartition.Error
	}

	log.Info().
		Str("topic", *m.TopicPartition.Topic).
		Int32("partition", m.TopicPartition.Partition).
		Msg("Kafka message delivered successfully")

	return nil
}
