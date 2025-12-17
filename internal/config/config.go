package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/semaphore"
)

var (
	PermissionGrpcHost string
	PostgresPoolSize   = 10
	PostgresHost       string
	PostgresPort       int
	PostgresUser       string
	PostgresPassword   string
	PostgresDb         string
	PostgresUserDB     string
	// Semaphore for Postgres queries
	PostgresSemaphore *semaphore.Weighted
)

// Redis Config
var (
	RedisHost     string
	RedisPort     int
	RedisUser     string
	RedisPassword string
	RedisDb       int
	RedisPoolSize = 10
	// Semaphore for Redis operations
	RedisSemaphore *semaphore.Weighted
)

var AdminEmails string

const RedisCachePrefix = "api.oms.cache:"
const DefaultCacheDuration = 5 * time.Second

// Kafka Config
var (
	KafkaServers         string
	KafkaConsumerGroup   = "oms.orderbook.consumer"
	KafkaMarketDataTopic = "market.data.transformed"
	KafkaMessageSource   = "dnse"
)

// Message types
const MessageTypeOrderBook = "TP"

func InitConfig() error {
	var err error

	// Load .env
	if err = godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found, using environment variables only")
	}

	// PERMISSION GRPC
	PermissionGrpcHost = os.Getenv("PERMISSION_GRPC_HOST")
	if PermissionGrpcHost == "" {
		return fmt.Errorf("PERMISSION_GRPC_HOST is not set")
	}

	// POSTGRES
	PostgresHost = os.Getenv("POSTGRES_HOST")
	PostgresPort, err = strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		return fmt.Errorf("failed to parse POSTGRES_PORT")
	}

	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	PostgresDb = os.Getenv("POSTGRES_DB")
	PostgresUserDB = os.Getenv("POSTGRES_USER_DB")
	// Init semaphore (max N concurrent Postgres queries)
	PostgresSemaphore = semaphore.NewWeighted(int64(PostgresPoolSize))

	// REDIS
	RedisHost = os.Getenv("REDIS_HOST")

	RedisPort, err = strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		return fmt.Errorf("failed to parse REDIS_PORT")
	}

	RedisUser = os.Getenv("REDIS_USER")
	RedisPassword = os.Getenv("REDIS_PASSWORD")

	RedisDb, err = strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return fmt.Errorf("failed to parse REDIS_DB")
	}

	// Init semaphore (max N concurrent Redis calls)
	RedisSemaphore = semaphore.NewWeighted(int64(RedisPoolSize))

	// Admin email
	AdminEmails = os.Getenv("ADMIN_EMAIL")

	// KAFKA
	KafkaServers = os.Getenv("KAFKA_SERVERS")
	if KafkaServers == "" {
		log.Warn().Msg("KAFKA_SERVERS not set, Kafka consumer will not be available")
	}

	return nil
}
