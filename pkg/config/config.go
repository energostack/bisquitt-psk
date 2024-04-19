package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog/log"
)

// Config represents the configuration of the application.
type Config struct {
	Host              string   `env:"HOST" envDefault:"localhost"`
	Port              string   `env:"PORT" envDefault:"3000"`
	FilePath          string   `env:"FILE_PATH" envDefault:"./data/psk.csv"`
	BasicAuthUsername string   `env:"BASIC_AUTH_USERNAME" envDefault:"admin"`
	BasicAuthPassword string   `env:"BASIC_AUTH_PASSWORD" envDefault:"admin"`
	KafkaTopic        string   `env:"KAFKA_TOPIC" envDefault:"clients"`
	KafkaBrokers      []string `env:"KAFKA_BROKERS" envDefault:"kafka:9092" envSeparator:","`
	KafkaGroup        string   `env:"KAFKA_GROUP" envDefault:"bisquitt_psk"`
	DataSource        string   `env:"DATA_SOURCE" required:"true"`
	APITimeoutSeconds int      `env:"API_TIMEOUT_SECONDS" envDefault:"5"`
	DocsEnabled       bool     `env:"DOCS_ENABLED" envDefault:"true"`
}

// New creates a new Config instance.
func New() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Err(err).Msg("Failed to parse config")
	}
	if len(cfg.DataSource) == 0 {
		log.Fatal().Msg("Data source is required")
	}
	return &cfg
}
