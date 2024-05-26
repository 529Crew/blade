package pump_monitor_hooks

import (
	"fmt"
	"strings"
	"time"

	"github.com/529Crew/blade/idls/pump"
	"github.com/529Crew/blade/internal/config"
	"github.com/529Crew/blade/internal/logger"
	"github.com/529Crew/blade/internal/requests"
	"github.com/529Crew/blade/internal/sol"
	pump_tx "github.com/529Crew/blade/internal/systems/pump/tx"
	"github.com/529Crew/blade/internal/types"
	"github.com/529Crew/blade/internal/util"
	"github.com/529Crew/blade/internal/webhooks"
	"github.com/gagliardetto/solana-go"
	"github.com/gtuk/discordwebhook"
)

func ParseCreateAndBuy(tx *solana.Transaction, sig string, preBalances []int64, postBalances []int64, postTokenBalances []types.TokenBalance) error {
	transaction, err := sol.ResolveAddressTables(tx)
	if err != nil {
		return err
	}
	tx = transaction

	var createInstruction *pump.Create
	parseCreateInst := func(inst *solana.CompiledInstruction) error {
		instAccs, err := inst.ResolveInstructionAccounts(&transaction.Message)
		if err != nil {
			return err
		}

		instruction, err := pump.DecodeInstruction(instAccs, inst.Data)
		if err != nil {
			return err
		}

		createInst, ok := instruction.Impl.(*pump.Create)
		if !ok {
			return fmt.Errorf("error casting instruction to create: %s", err)
		}
		createInstruction = createInst

		return nil
	}

	var buyInstruction *pump.Buy
	parseBuyInst := func(inst *solana.CompiledInstruction) error {
		instAccs, err := inst.ResolveInstructionAccounts(&transaction.Message)
		if err != nil {
			return err
		}

		instruction, err := pump.DecodeInstruction(instAccs, inst.Data)
		if err != nil {
			return err
		}

		buyInst, ok := instruction.Impl.(*pump.Buy)
		if !ok {
			return fmt.Errorf("error casting instruction to buy: %s", err)
		}
		buyInstruction = buyInst

		return nil
	}

	/* search for create or buy inst */
	for _, inst := range tx.Message.Instructions {
		if len([]byte(inst.Data)) < 8 {
			continue
		}
		var discriminator [8]byte
		copy(discriminator[:], []byte(inst.Data)[:8])

		switch {
		case discriminator == [8]byte{24, 30, 200, 40, 5, 28, 7, 119}: /* create */
			parseCreateInst(&inst)
		case discriminator == [8]byte{102, 6, 61, 18, 1, 218, 235, 234}: /* buy */
			parseBuyInst(&inst)
		}
	}

	/* ensure both instructions were parsed */
	if createInstruction == nil || buyInstruction == nil {
		return fmt.Errorf("failed to parse create or buy inst")
	}

	/* get ipfs metadata */
	metadata, err := requests.IpfsData(*createInstruction.Uri)
	if err != nil {
		return fmt.Errorf("error getting pf data: %s", err)
	}

	/* get creator's coins */
	coins, err := requests.Coins(createInstruction.GetUserAccount().PublicKey.String())
	if err != nil {
		return fmt.Errorf("error getting coins: %s", err)
	}

	/* get pre / post sol balance */
	preSolBalance := float64(preBalances[0]) / 1_000_000_000
	postSolBalance := float64(postBalances[0]) / 1_000_000_000

	/* get token balance */
	var tokenBalance float64 = 0
	for _, postTokenBalance := range postTokenBalances {
		if postTokenBalance.Owner == createInstruction.GetUserAccount().PublicKey.String() && postTokenBalance.Mint == createInstruction.GetMintAccount().PublicKey.String() {
			tokenBalance = postTokenBalance.UITokenAmount.UIAmount
		}
	}

	/* dev coin metrics */
	totalTokens := 0
	totalKoth := 0
	totalRaydium := 0
	seenMints := make(map[string]bool)
	foundCurrent := false

	for _, coin := range *coins {
		if coin.Mint == createInstruction.GetMintAccount().PublicKey.String() {
			foundCurrent = true
		}

		/* check if we already saw this mint */
		_, ok := seenMints[coin.Mint]
		if ok {
			continue
		}
		seenMints[coin.Mint] = true

		totalTokens++

		if coin.KingOfTheHillTimestamp != 0 {
			totalKoth++
		}

		if coin.RaydiumPool != "" {
			totalRaydium++
		}
	}
	if !foundCurrent {
		totalTokens++
	}

	/* pre-calculate more metrics */
	solSpent := preSolBalance - postSolBalance
	percentOwned := (tokenBalance / 1_000_000_000) * 100

	go sendCreateAndBuyWebhook(createInstruction, sig, metadata, postSolBalance, tokenBalance, solSpent, percentOwned, totalTokens, totalKoth, totalRaydium, config.Get().PfCreateWebhook)
	processCreateAndBuy(createInstruction, buyInstruction, sig, metadata, postSolBalance, tokenBalance, solSpent, percentOwned, totalTokens, totalKoth, totalRaydium)

	return nil
}

// var boughtAlready bool = false

func processCreateAndBuy(
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
	mint := createInst.GetMintAccount().PublicKey.String()

	/* check koth and raydium stats */
	if totalTokens > cfg.MaximumTotal {
		logger.Log.Printf("[FILTERED]: mint %s - total tokens over max", mint)
		return
	}
	if totalKoth < cfg.MinimumKoth {
		logger.Log.Printf("[FILTERED]: mint %s - koth under min", mint)
		return
	}
	if totalRaydium < cfg.MinimumRaydium {
		logger.Log.Printf("[FILTERED]: mint %s - ray under min", mint)
		return
	}

	/* check dev info */
	if postSolBalance < cfg.DevMinimumSolBalance {
		logger.Log.Printf("[FILTERED]: mint %s - dev bal under min", mint)
		return
	}
	if percentOwned > cfg.DevMaximumPercent {
		logger.Log.Printf("[FILTERED]: mint %s - dev percent over max", mint)
		return
	}

	/* check socials */
	if cfg.WebsiteRequired && metadata.Website == "" {
		logger.Log.Printf("[FILTERED]: mint %s - no website", mint)
		return
	}
	if cfg.TwitterRequired && metadata.Twitter == "" {
		logger.Log.Printf("[FILTERED]: mint %s - no twitter", mint)
		return
	}
	if cfg.TelegramRequired && metadata.Telegram == "" {
		logger.Log.Printf("[FILTERED]: mint %s - no telegram", mint)
		return
	}

	/* check banned words in title / description */
	for _, bannedWord := range cfg.BannedWords {
		if strings.Contains(metadata.Name, bannedWord) {
			logger.Log.Printf("[FILTERED]: mint %s - banned word found in name / %s", mint, bannedWord)
			return
		}
		if strings.Contains(metadata.Description, bannedWord) {
			logger.Log.Printf("[FILTERED]: mint %s - banned word found in description / %s", mint, bannedWord)
			return
		}
	}

	go sendCreateAndBuyWebhook(createInst, sig, metadata, postSolBalance, tokenBalance, solSpent, percentOwned, totalTokens, totalKoth, totalRaydium, config.Get().PfFilteredCreateWebhook)

	// if !boughtAlready {
	// 	boughtAlready = true
	pump_tx.Buy(buyInst, solSpent, tokenBalance)
	// }
}

func sendCreateAndBuyWebhook(
	inst *pump.Create,
	sig string,
	metadata *types.IpfsResponse,
	postSolBalance float64,
	tokenBalance float64,
	solSpent float64,
	percentOwned float64,
	totalTokens int,
	totalKoth int,
	totalRaydium int,
	webhook string,
) {
	fields := []discordwebhook.Field{
		{
			Name:  webhooks.StrPtr("Name"),
			Value: webhooks.StrPtr(fmt.Sprintf("```%s```", metadata.Name)),
		},
		{
			Name:  webhooks.StrPtr("Ticker"),
			Value: webhooks.StrPtr(fmt.Sprintf("```%s```", metadata.Symbol)),
		},
	}

	/* add description if not empty */
	if len(metadata.Description) > 1000 {
		metadata.Description = metadata.Description[:1000]
	}
	if metadata.Description != "" {
		fields = append(fields, discordwebhook.Field{
			Name:  webhooks.StrPtr("Description"),
			Value: webhooks.StrPtr(fmt.Sprintf("```%s```", metadata.Description)),
		})
	}

	/* add basic token fields */
	fields = append(fields, []discordwebhook.Field{
		{
			Name:  webhooks.StrPtr("Token Address"),
			Value: webhooks.StrPtr(fmt.Sprintf("```%s```", inst.GetMintAccount().PublicKey.String())),
		},
	}...)

	/* add dev coin stats */
	fields = append(fields, []discordwebhook.Field{
		{
			Name:  webhooks.StrPtr("Dev Address"),
			Value: webhooks.StrPtr(fmt.Sprintf("```%s```", inst.GetUserAccount().PublicKey.String())),
		},
		{
			Name:  webhooks.StrPtr("Dev Coin Stats"),
			Value: webhooks.StrPtr(fmt.Sprintf("```Total %d / KOTH %d / Raydium %d```", totalTokens, totalKoth, totalRaydium)),
		},
	}...)

	/* add dev balances */
	fields = append(fields, []discordwebhook.Field{
		{
			Name:  webhooks.StrPtr("Dev Balances"),
			Value: webhooks.StrPtr(fmt.Sprintf("```Spent         | %.2f SOL\nBalance       | %.2f SOL\nToken Balance | %.2f (%.2f%% Of Supply)\n```", solSpent, postSolBalance, tokenBalance, percentOwned)),
		},
	}...)

	/* add socials if available */
	var socialsStr []string
	socialsStr = append(socialsStr, fmt.Sprintf("[PUMP FUN](https://pump.fun/%s)", inst.GetMintAccount().PublicKey.String()))
	socialsStr = append(socialsStr, fmt.Sprintf("[PHOTON](https://photon-sol.tinyastro.io/en/lp/%s)", inst.GetMintAccount().PublicKey.String()))
	if metadata.Telegram != "" {
		if !strings.HasPrefix(metadata.Telegram, "http") {
			metadata.Telegram = fmt.Sprintf("https://%s", metadata.Telegram)
		}
		socialsStr = append(socialsStr, fmt.Sprintf("[TELEGRAM](%s)", metadata.Telegram))
	}
	if metadata.Twitter != "" {
		if !strings.HasPrefix(metadata.Twitter, "http") {
			metadata.Twitter = fmt.Sprintf("https://%s", metadata.Twitter)
		}
		socialsStr = append(socialsStr, fmt.Sprintf("[TWITTER](%s)", metadata.Twitter))
	}
	if metadata.Website != "" {
		if !strings.HasPrefix(metadata.Website, "http") {
			metadata.Website = fmt.Sprintf("https://%s", metadata.Website)
		}
		socialsStr = append(socialsStr, fmt.Sprintf("[WEBSITE](%s)", metadata.Website))
	}
	if len(socialsStr) > 0 {
		fields = append(fields, discordwebhook.Field{
			Name:  webhooks.StrPtr("Links"),
			Value: webhooks.StrPtr(strings.Join(socialsStr, " | ")),
		})
	}

	/* add image if available */
	var thumbnail *discordwebhook.Thumbnail = nil
	if metadata.Image != "" {
		if strings.Contains(metadata.Image, "/ipfs/") {
			urlSplit := strings.Split(metadata.Image, "/ipfs/")
			if len(urlSplit) > 1 {
				thumbnail = &discordwebhook.Thumbnail{
					Url: webhooks.StrPtr(fmt.Sprintf("https://flowgocrazy.mypinata.cloud/ipfs/%s", urlSplit[1])),
				}
			}
		} else {
			thumbnail = &discordwebhook.Thumbnail{
				Url: webhooks.StrPtr(metadata.Image),
			}
		}
	}

	message := discordwebhook.Message{
		Username:  webhooks.StrPtr("After Hours Monitors"),
		AvatarUrl: webhooks.StrPtr(webhooks.AvatarURL),
		Embeds: &[]discordwebhook.Embed{
			{
				Title: webhooks.StrPtr("New Pump Fun Token"),
				Url:   webhooks.StrPtr(fmt.Sprintf("https://solscan.io/tx/%s", sig)),
				Color: webhooks.StrPtr("2303786"),

				Fields: &fields,

				Thumbnail: thumbnail,

				Footer: &discordwebhook.Footer{
					Text:    webhooks.StrPtr(fmt.Sprintf("After Hours Monitors - %s", time.Now().UTC().In(logger.TimeLocation).Format("Mon Jan 2 03:04:05 PM EST"))),
					IconUrl: webhooks.StrPtr(webhooks.AvatarURL),
				},
			},
		},
	}

	cfg := config.Get()

	if cfg.WebhooksEnabled {
		err := discordwebhook.SendMessage(webhook, message)
		if err != nil {
			if !strings.Contains(err.Error(), "rate limited") {
				logger.Log.Println(err)
				util.PrettyPrint(message)
			}
		}
	}
}
