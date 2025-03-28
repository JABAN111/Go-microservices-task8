package config

type Config struct {
	BindAddress string `yaml:"words_address" env:"WORDS_ADDRESS" default:"localhost:8080"`
}

func NewConfig() *Config {
	return &Config{
		BindAddress: "",
	}
}

func DefaultConfig() *Config {
	return &Config{
		BindAddress: "localhost:8080",
	}
}
