package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Services struct {
	WordsAddress  string `yaml:"words_address" env:"WORDS_ADDRESS" env-default:":8080"`
	UpdateAddress string `json:"update_address" env:"UPDATE_ADDRESS" env-default:":8080"`
}

type Config struct {
	Services  `yaml:"services"`
	DBAddress string `yaml:"db_address" env:"DB_ADDRESS"`
	Workers   int    `yaml:"concurrency" env:"CONCURRENCY" env-default:"30"`
	Address   string `yaml:"search_address" env:"SEARCH_ADDRESS" env-default:":8080"`
	LogLevel  string `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
}

func MustLoad(configPath string) Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config %q: %s", configPath, err)
	}
	return cfg
}
