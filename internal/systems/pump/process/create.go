package pump_process

import (
	"strings"

	"github.com/529Crew/blade/idls/pump"
	"github.com/529Crew/blade/internal/config"
	pump_tx "github.com/529Crew/blade/internal/systems/pump/tx"
	"github.com/529Crew/blade/internal/types"
)

func ProcessCreateAndBuy(
	createInst *pump.Create,
	buyInst *pump.Buy,
	sig string,
	metadata *types.IpfsResponse,
	postSolBalance float64,
	tokenBalance float64,
	solSpent float64,
	percentOwned float64,
	totalTokens int,
	totalKoth int,
	totalRaydium int,
) {
	cfg := config.Get()

	/* check koth and raydium stats */
	if totalTokens > cfg.MaximumTotal {
		return
	}
	if totalKoth < cfg.MinimumKoth {
		return
	}
	if totalRaydium < cfg.MinimumRaydium {
		return
	}

	/* check dev info */
	if postSolBalance < cfg.DevMinimumSolBalance {
		return
	}
	if percentOwned > cfg.DevMaximumPercent {
		return
	}

	/* check socials */
	if cfg.WebsiteRequired && metadata.Website == "" {
		return
	}
	if cfg.TwitterRequired && metadata.Twitter == "" {
		return
	}
	if cfg.TelegramRequired && metadata.Telegram == "" {
		return
	}

	/* check banned words in title / description */
	for _, bannedWord := range cfg.BannedWords {
		if strings.Contains(metadata.Name, bannedWord) {
			return
		}
		if strings.Contains(metadata.Description, bannedWord) {
			return
		}
	}

	pump_tx.Buy(buyInst)
}
