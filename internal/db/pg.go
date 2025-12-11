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

func InitPostgres() error {
	Postgres = pg.Connect(&pg.Options{
		Addr:         fmt.Sprintf("%s:%d", config.PostgresHost, config.PostgresPort),
		User:         config.PostgresUser,
		Password:     config.PostgresPassword,
		Database:     config.PostgresDb,
		PoolSize:     config.PostgresPoolSize,
		IdleTimeout:  5 * time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	err := Postgres.Ping(context.Background())
	if err != nil {
		log.Error().Msgf("Error connecting to Postgres: %v", err)
		return err
	} else {
		log.Info().Msgf("Connected to Postgres at %s:%d/{%s}", config.PostgresHost, config.PostgresPort, config.PostgresDb)
		return nil
	}
}

// ClosePostgres
func ClosePostgres() {
	err := Postgres.Close()
	if err != nil {
		return
	}
	log.Info().Msgf("Closed Postgres")
}
