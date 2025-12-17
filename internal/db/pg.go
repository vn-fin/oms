package db

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog/log"
	"github.com/vn-fin/oms/internal/config"
)

var Postgres *pg.DB
var PostgresUserDB *pg.DB

func InitPostgres() error {
	Postgres = pg.Connect(&pg.Options{
		Addr:         fmt.Sprintf("%s:%d", config.PostgresHost, config.PostgresPort),
		User:         config.PostgresUser,
		Password:     config.PostgresPassword,
		Database:     config.PostgresDb,
		PoolSize:     config.PostgresPoolSize,
		MaxRetries:   5,
		IdleTimeout:  5 * time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	err := Postgres.Ping(context.Background())
	if err != nil {
		log.Error().Msgf("Error connecting to Postgres: %v", err)
		return err
	}
	log.Info().Msgf("Connected to Postgres at %s:%d/{%s}", config.PostgresHost, config.PostgresPort, config.PostgresDb)
	return nil
}

func InitPostgresUserDB() error {
	PostgresUserDB = pg.Connect(&pg.Options{
		Addr:         fmt.Sprintf("%s:%d", config.PostgresHost, config.PostgresPort),
		User:         config.PostgresUser,
		Password:     config.PostgresPassword,
		Database:     config.PostgresUserDB,
		PoolSize:     config.PostgresPoolSize,
		MaxRetries:   5,
		IdleTimeout:  5 * time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	err := PostgresUserDB.Ping(context.Background())
	if err != nil {
		log.Error().Msgf("Error connecting to Postgres User DB: %v", err)
		return err
	}
	log.Info().Msgf("Connected to Postgres User DB at %s:%d/{%s}", config.PostgresHost, config.PostgresPort, config.PostgresUserDB)
	return nil
}

// ClosePostgres
func ClosePostgres() {
	if Postgres != nil {
		err := Postgres.Close()
		if err != nil {
			log.Error().Msgf("Error closing Postgres: %v", err)
			return
		}
		log.Info().Msgf("Closed Postgres")
	}

	if PostgresUserDB != nil {
		err := PostgresUserDB.Close()
		if err != nil {
			log.Error().Msgf("Error closing Postgres User DB: %v", err)
			return
		}
		log.Info().Msgf("Closed Postgres User DB")
	}
}
