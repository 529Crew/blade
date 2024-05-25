package config

import (
	"os"

	"github.com/529Crew/blade/internal/logger"
	"github.com/BurntSushi/toml"
)

type Config struct {
	/* basic bot config */

	HeliusApiKey  string `toml:"helius_api_key"`
	GeyserGrpcUrl string `toml:"geyser_grpc_url"`

	UtilRpcUrl string `toml:"util_rpc_url"`
	RpcUrl     string `toml:"rpc_url"`
	JitoRpcUrl string `toml:"jito_rpc_url"`

	IpfsGatewayUri string `toml:"ipfs_gateway_uri"`

	JitoPrivateKey   string `toml:"jito_private_key"`
	WalletPrivateKey string `toml:"wallet_private_key"`

	BuyPriorityFee uint64 `toml:"buy_priority_fee"`

	WebhooksEnabled bool `toml:"webhooks_enabled"`

	/* webhooks */

	SelfMonitorWebhook      string `toml:"self_monitor_webhook"`
	PfCreateWebhook         string `toml:"pf_create_webhook"`
	PfFilteredCreateWebhook string `toml:"pf_filtered_create_webhook"`

	/* token filters */

	MaximumTotal   int `toml:"maximum_total"`
	MinimumKoth    int `toml:"minimum_koth"`
	MinimumRaydium int `toml:"minimum_raydium"`

	DevMinimumSolBalance float64 `toml:"dev_minimum_sol_balance"`
	DevMaximumPercent    float64 `toml:"dev_maximum_percent"`

	WebsiteRequired  bool `toml:"website_required"`
	TwitterRequired  bool `toml:"twitter_required"`
	TelegramRequired bool `toml:"telegram_required"`

	BannedWords []string `toml:"banned_words"`

	/* tx settings */

	BuyAmount           uint64 `toml:"buy_amount"`
	BuySlippage         int64  `toml:"buy_slippage"`
	BloxrouteTip        uint64 `toml:"bloxroute_tip"`
	BloxrouteAuthHeader string `toml:"bloxroute_auth_header"`
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
