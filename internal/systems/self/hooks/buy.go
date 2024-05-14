package self_monitor_hooks

import (
	"fmt"
	"strings"
	"time"

	"github.com/529Crew/blade/idls/pump"
	"github.com/529Crew/blade/internal/config"
	helius_api "github.com/529Crew/blade/internal/helius/api"
	"github.com/529Crew/blade/internal/logger"
	"github.com/529Crew/blade/internal/requests"
	"github.com/529Crew/blade/internal/types"
	"github.com/529Crew/blade/internal/util"
	"github.com/529Crew/blade/internal/webhooks"
	"github.com/gagliardetto/solana-go"
	"github.com/gtuk/discordwebhook"
)

func ParseBuy(tx *solana.Transaction, sig string, preBalances []int64, postBalances []int64, preTokenBalances []types.TokenBalance, postTokenBalances []types.TokenBalance) error {

	parseBuyInst := func(inst *solana.CompiledInstruction) error {
		instAccs, err := inst.ResolveInstructionAccounts(&tx.Message)
		if err != nil {
			return err
		}

		instruction, err := pump.DecodeInstruction(instAccs, inst.Data)
		if err != nil {
			return err
		}

		buyInst, ok := instruction.Impl.(*pump.Buy)
		if !ok {
			return fmt.Errorf("error casting instruction to buy: %v", err)
		}

		/* get asset metadata ipfs url */
		mplMetadata, err := helius_api.GetAsset(buyInst.GetMintAccount().PublicKey.String())
		if err != nil {
			return fmt.Errorf("error getting asset: %v", err)
		}

		/* get ipfs metadata */
		metadata, err := requests.IpfsData(mplMetadata.Result.Content.JSONURI)
		if err != nil {
			return fmt.Errorf("error getting pf data: %v", err)
		}

		/* calculate sol spent */
		preSolBalance := preBalances[0]
		postSolBalance := postBalances[0]
		solSpent := float64(preSolBalance-postSolBalance) / 1_000_000_000

		/* calculate tokens received */
		var preTokenBalance float64 = 0
		for _, balance := range preTokenBalances {
			if balance.Mint == buyInst.GetMintAccount().PublicKey.String() && balance.Owner == buyInst.GetUserAccount().PublicKey.String() {
				preTokenBalance = balance.UITokenAmount.UIAmount
			}
		}
		var postTokenBalance float64 = 0
		for _, balance := range postTokenBalances {
			if balance.Mint == buyInst.GetMintAccount().PublicKey.String() && balance.Owner == buyInst.GetUserAccount().PublicKey.String() {
				postTokenBalance = balance.UITokenAmount.UIAmount
			}
		}
		tokensReceived := postTokenBalance - preTokenBalance

		sendBuyWebhook(buyInst, sig, mplMetadata, metadata, solSpent, tokensReceived)

		return nil
	}

	/* search for buy inst */
	for _, inst := range tx.Message.Instructions {
		if len([]byte(inst.Data)) < 8 {
			continue
		}
		var discriminator [8]byte
		copy(discriminator[:], []byte(inst.Data)[:8])

		/* if discriminator is buy */
		if discriminator == [8]byte{102, 6, 61, 18, 1, 218, 235, 234} {
			return parseBuyInst(&inst)
		}
	}

	return nil
}

var SELF_BUY_WEBHOOK = "https://discord.com/api/webhooks/1239965126237360221/cgS5LqfyaQnMd7hnaNE6vBBRM7uv704QkufJRRUfW9N_e5jJh8P6Bru3nD7GVb-sYS-F"

func sendBuyWebhook(inst *pump.Buy, sig string, mplMetadata *types.GetAssetResponse, metadata *types.IpfsResponse, solSpent float64, tokensReceived float64) {
	fields := []discordwebhook.Field{
		{
			Name:  webhooks.StrPtr("Name"),
			Value: webhooks.StrPtr(fmt.Sprintf("```%s```", mplMetadata.Result.Content.Metadata.Name)),
		},
		{
			Name:  webhooks.StrPtr("Ticker"),
			Value: webhooks.StrPtr(fmt.Sprintf("```%s```", mplMetadata.Result.Content.Metadata.Symbol)),
		},
	}

	/* add description if not empty */
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

	/* add buy info */
	fields = append(fields, []discordwebhook.Field{
		{
			Name:  webhooks.StrPtr("SOL Spent"),
			Value: webhooks.StrPtr(fmt.Sprintf("```%.2f SOL```", solSpent)),
		},
		{
			Name:  webhooks.StrPtr("Tokens Received"),
			Value: webhooks.StrPtr(fmt.Sprintf("```%.2f %s```", tokensReceived, mplMetadata.Result.Content.Metadata.Symbol)),
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
				Title: webhooks.StrPtr("Bought Tokens"),
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

	if config.Get().WebhooksEnabled {
		err := discordwebhook.SendMessage(SELF_BUY_WEBHOOK, message)
		if err != nil {
			if !strings.Contains(err.Error(), "rate limited") {
				logger.Log.Println(err)
				util.PrettyPrint(message)
			}
		}
	}
}
