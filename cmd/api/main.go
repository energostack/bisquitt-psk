package main

import (
	"context"
	"os"
	"os/signal"
	"runtime/debug"

	"github.com/energostack/bisquitt-psk/pkg/api"
	"github.com/energostack/bisquitt-psk/pkg/config"
	"github.com/energostack/bisquitt-psk/pkg/mapstore"
	"github.com/energostack/bisquitt-psk/pkg/reader"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//	@title						Bisquitt PSK API
//	@version					1.0
//	@description				This is OpenAPI(2.0) for Bisquitt PSK
//	@securityDefinitions.basic	BasicAuth
//	@host						localhost:3000
//	@BasePath					/

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Info().Msgf("Recovered from panic: %v", r)
			log.Info().Msgf("Stack trace: %s", debug.Stack())
		}
	}()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	cfg := config.New()

	var dataReader reader.Reader
	switch cfg.DataSource {
	case "csv":
		log.Info().Msg("Using CSV as data source")
		csvReader, err := reader.NewCSVReader(cfg.FilePath, ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create CSV reader")
			return
		}
		dataReader = csvReader
	case "sqlite":
		log.Info().Msg("Using SQLite as data source")
		sqliteReader, err := reader.NewSQLiteReader(cfg.FilePath, ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create SQLite reader")
			return
		}
		dataReader = sqliteReader
	case "kafka":
		log.Info().Msg("Using Kafka as data source")
		kafkaReader, err := reader.NewKafkaReader(cfg, ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create Kafka reader")
			return
		}
		dataReader = kafkaReader
	default:
		log.Fatal().Msg("Unknown data source")
	}

	controller := mapstore.NewController(dataReader)
	handler := api.NewCustomHandler(controller)

	r := api.NewCustomRouter(handler, cfg)

	err := r.ListenAndServe(cfg.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start HTTP server")
		return
	}
}
