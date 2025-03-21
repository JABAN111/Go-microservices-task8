package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServerConfig struct {
	ServerAddress string        `yaml:"address" env:"API_ADDRESS"`
	HttpTimeout   time.Duration `yaml:"timeout" env:"HTTP_SERVER_TIMEOUT"`
}

type Config struct {
	HttpServer    HttpServerConfig `yaml:"http_server"`
	WordsAddress  string           `yaml:"words_address" env:"WORDS_ADDRESS"`
	UpdateAddress string           `json:"update_address" env:"UPDATE_ADDRESS"`
	LogLevel      string           `yaml:"log_level" env:"LOG_LEVEL"`
}

func defaultConfig() Config {
	return Config{
		HttpServer: HttpServerConfig{
			ServerAddress: "0.0.0.0:28080",
			HttpTimeout:   time.Hour / 2,
		},
		WordsAddress:  "localhost:28081",
		LogLevel:      "info",
		UpdateAddress: "localhost:28082",
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
