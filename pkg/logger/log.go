package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupLogger() {
	// Set the global minimum log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Add file and line number automatically
	log.Logger = log.With().Caller().Logger()
}
