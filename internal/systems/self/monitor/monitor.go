package self_monitor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/529Crew/blade/internal/logger"
	self_system "github.com/529Crew/blade/internal/systems/self"
	self_monitor_hooks "github.com/529Crew/blade/internal/systems/self/hooks"
	"github.com/529Crew/blade/internal/types"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"

	helius_ws "github.com/529Crew/blade/internal/helius/ws"
)

func Monitor(connected chan struct{}) {
	_, listener := helius_ws.Listen()

	var subId int64

	for {
		select {
		case <-connected:
			txSubId, err := txSubscribe()
			if err != nil {
				log.Printf("[SELF MONITOR]: txsub error: %v\n", err)
			}
			subId = txSubId
			logger.Log.Println("[SELF MONITOR]: activated")
		case msg := <-listener:
			/* filter all msg without matching txSubId */
			if !strings.Contains(string(msg), fmt.Sprintf(`"subscription":%d`, subId)) {
				continue
			}
			go sort(msg)
		}
	}
}

func txSubscribe() (txSubId int64, e error) {
	return helius_ws.TransactionSubscribe(&types.TransactionSubscribePayload{
		Jsonrpc: "2.0",
		Method:  "transactionSubscribe",
		Params: []types.TransactionSubscribeParams{
			{
				Vote:   false,
				Failed: false,
				AccountRequired: []string{
					self_system.SELF.String(),
				},
			},
			{
				Commitment:                     "processed",
				Encoding:                       "base64",
				TransactionDetails:             "full",
				ShowRewards:                    false,
				MaxSupportedTransactionVersion: 0,
			},
		},
	})
}

func sort(msg []byte) {
	var notification types.TransactionNotification
	err := json.Unmarshal(msg, &notification)
	if err != nil {
		log.Printf("[SELF MONITOR]: error unmarshalling transaction notification: %v\n", err)
		return
	}

	combinedLogs := strings.Join(notification.Params.Result.Transaction.Meta.LogMessages, " ")
	sig := notification.Params.Result.Signature

	var txType string
	switch {
	case strings.Contains(combinedLogs, "Instruction: Buy"):
		txType = "buy"
	case strings.Contains(combinedLogs, "Instruction: Sell"):
		txType = "sell"
	}

	if txType == "" {
		logger.Log.Printf("[SELF MONITOR]: %s", sig)
		return
	}

	data, err := base64.StdEncoding.DecodeString(notification.Params.Result.Transaction.Transaction[0])
	if err != nil {
		logger.Log.Printf("[SELF MONITOR]: %s - %s / %s", txType, sig, err)
		return
	}

	tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(data))
	if err != nil {
		logger.Log.Printf("[SELF MONITOR]: %s - %s / %s", txType, sig, err)
		return
	}
	preBalances, postBalances := notification.Params.Result.Transaction.Meta.PreBalances, notification.Params.Result.Transaction.Meta.PostBalances
	preTokenBalances, postTokenBalances := notification.Params.Result.Transaction.Meta.PreTokenBalances, notification.Params.Result.Transaction.Meta.PostTokenBalances

	switch {
	case txType == "buy":
		err = self_monitor_hooks.ParseBuy(tx, sig, preBalances, postBalances, preTokenBalances, postTokenBalances)
	case txType == "sell":
		err = self_monitor_hooks.ParseSell(tx, sig, preBalances, postBalances, preTokenBalances, postTokenBalances)
	}
	if err != nil {
		logger.Log.Printf("[SELF MONITOR]: %s - %s / %s", txType, sig, err)
	}
}
