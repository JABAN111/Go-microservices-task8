package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServerConfig struct {
	ServerAddress string        `yaml:"address" env:"API_ADDRESS"`
	HttpTimeout   time.Duration `yaml:"timeout" env:"HTTP_SERVER_TIMEOUT" env-default:"20s"`
}

type Config struct {
	HttpServer         HttpServerConfig `yaml:"http_server"`
	WordsAddress       string           `yaml:"words_address" env:"WORDS_ADDRESS"`
	UpdateAddress      string           `json:"update_address" env:"UPDATE_ADDRESS"`
	SearchAddress      string           `json:"search_address" env:"SEARCH_ADDRESS"`
	ConcurrencyLimiter int              `json:"concurrency_limiter" env:"SEARCH_CONCURRENCY" env-default:"10"`
	RateLimiter        int              `json:"rate_limiter" env:"SEARCH_RATE" env-default:"100"`
	LogLevel           string           `yaml:"log_level" env:"LOG_LEVEL"`
	TokenTTL           time.Duration    `yaml:"token_ttl" env:"TOKEN_TTL" env-default:"2m"`
}

func defaultConfig() Config {
	return Config{
		HttpServer: HttpServerConfig{
			ServerAddress: "0.0.0.0:28080",
			HttpTimeout:   time.Hour / 2,
		},
		WordsAddress:  "localhost:28081",
		LogLevel:      "DEBUG",
		UpdateAddress: "localhost:28082",
		SearchAddress: "localhost:28087",
	}
}

func MustLoad(configPath string) Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Println("[WARN] error parse config file, using default configuration")
		return defaultConfig()
	}
	return cfg
}
