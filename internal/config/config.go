package config

import (
	"os"

	"github.com/529Crew/blade/internal/logger"
	"github.com/BurntSushi/toml"
)

type Config struct {
	HeliusApiKey  string `toml:"helius_api_key"`
	GeyserGrpcUrl string `toml:"geyser_grpc_url"`

	UtilRpcUrl string `toml:"util_rpc_url"`
	RpcUrl     string `toml:"rpc_url"`

	JitoPrivateKey string `toml:"jito_private_key"`
	JitoRpcUrl     string `toml:"jito_rpc_url"`

	IpfsGatewayUri string `toml:"ipfs_gateway_uri"`

	WebhooksEnabled bool `toml:"webhooks_enabled"`
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
