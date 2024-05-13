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
	"github.com/529Crew/blade/internal/types"
	"github.com/529Crew/blade/internal/util"
	"github.com/529Crew/blade/internal/webhooks"
	"github.com/gagliardetto/solana-go"
	"github.com/gtuk/discordwebhook"
)

func ParseCreateAndBuy(tx *solana.Transaction, sig string) error {
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
			return fmt.Errorf("error casting instruction to create: %v", err)
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
			return fmt.Errorf("error casting instruction to buy: %v", err)
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
		return fmt.Errorf("error getting pf data: %v", err)
	}

	sendCreateAndBuyWebhook(createInstruction, buyInstruction, sig, metadata)

	return nil
}

var PF_CREATE_AND_BUY_WEBHOOK = "https://discord.com/api/webhooks/1239698725174120458/DiLcFDxGIrZMXfk2nOfyN4INlS-5jH5JG0igmoKsqNweKIz_2z0_SlNCooKoVqXzenjj"

func sendCreateAndBuyWebhook(create *pump.Create, buy *pump.Buy, sig string, metadata *types.IpfsResponse) {
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
			Value: webhooks.StrPtr(fmt.Sprintf("```%s```", create.GetMintAccount().PublicKey.String())),
		},
	}...)

	/* add socials if available */
	var socialsStr []string
	socialsStr = append(socialsStr, fmt.Sprintf("[PUMP FUN](https://pump.fun/%s)", create.GetMintAccount().PublicKey.String()))
	socialsStr = append(socialsStr, fmt.Sprintf("[PHOTON](https://photon-sol.tinyastro.io/en/lp/%s)", create.GetMintAccount().PublicKey.String()))
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
		urlSplit := strings.Split(metadata.Image, "/ipfs/")
		if len(urlSplit) > 1 {
			thumbnail = &discordwebhook.Thumbnail{
				Url: webhooks.StrPtr(fmt.Sprintf("https://flowgocrazy.mypinata.cloud/ipfs/%s", urlSplit[1])),
			}
		}
	}

	message := discordwebhook.Message{
		Username:  webhooks.StrPtr("529 Monitors"),
		AvatarUrl: webhooks.StrPtr(webhooks.AvatarURL),
		Embeds: &[]discordwebhook.Embed{
			{
				Title: webhooks.StrPtr("Pump Fun LP Live On Raydium"),
				Url:   webhooks.StrPtr(fmt.Sprintf("https://solscan.io/tx/%s", sig)),
				Color: webhooks.StrPtr("2303786"),

				Fields: &fields,

				Thumbnail: thumbnail,

				Footer: &discordwebhook.Footer{
					Text:    webhooks.StrPtr(fmt.Sprintf("529 Monitors - %s", time.Now().UTC().In(logger.TimeLocation).Format("Mon Jan_2 03:04:05 PM EST"))),
					IconUrl: webhooks.StrPtr(webhooks.AvatarURL),
				},
			},
		},
	}

	if config.Get().WebhooksEnabled {
		err := discordwebhook.SendMessage(PF_CREATE_AND_BUY_WEBHOOK, message)
		if err != nil {
			if !strings.Contains(err.Error(), "rate limited") {
				logger.Log.Println(err)
				util.PrettyPrint(message)
			}
		}
	}
}
