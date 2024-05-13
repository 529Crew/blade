package config

import (
	"os"

	"github.com/529Crew/blade/internal/logger"
	"github.com/BurntSushi/toml"
)

type Config struct {
	HeliusApiKey  string `toml:"helius_api_key"`
	GeyserGrpcUrl string `toml:"geyser_grpc_url"`
}

var cfg *Config

func init() {
	logger.Log.Println("loading config from data/config.toml")

	bytes, err := os.ReadFile("data/config.toml")
	if err != nil {
		logger.Log.Panicf("failed to read data/config.toml: %s", err)
	}

	var config Config
	_, err = toml.Decode(string(bytes), &config)
	if err != nil {
		logger.Log.Panicf("failed to decode data/config.toml: %s", err)
	}
	cfg = &config

	logger.Log.Println("loaded config from data/config.toml")
}

func Get() *Config {
	return cfg
}
