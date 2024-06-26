package pump_monitor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/529Crew/blade/internal/logger"
	pump_monitor_hooks "github.com/529Crew/blade/internal/systems/pump/hooks"
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
				log.Printf("[PUMP MONITOR HELIUS]: txsub error: %v\n", err)
			}
			subId = txSubId
			logger.Log.Println("[PUMP MONITOR HELIUS]: activated")
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
					"6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P",
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
		log.Printf("[PUMP MONITOR HELIUS]: error unmarshalling transaction notification: %v\n", err)
		return
	}

	combinedLogs := strings.Join(notification.Params.Result.Transaction.Meta.LogMessages, " ")
	sig := notification.Params.Result.Signature

	/* check if we've already seen this sig */
	seen := Seen(sig)
	if seen {
		return
	}

	var txType string
	switch {
	case strings.Contains(combinedLogs, "Instruction: Create") && strings.Contains(combinedLogs, "Instruction: Buy"):
		txType = "create + buy"
	case strings.Contains(combinedLogs, "Instruction: Create") && !strings.Contains(combinedLogs, "Instruction: Buy"):
		txType = "create"
	case !strings.Contains(combinedLogs, "Instruction: Create") && strings.Contains(combinedLogs, "Instruction: Buy"):
		txType = "buy"
	case strings.Contains(combinedLogs, "Instruction: Sell"):
		txType = "sell"
	}

	if txType != "create + buy" {
		return
	}

	logger.Log.Printf("[PUMP MONITOR HELIUS]: detected %s / %s", txType, sig)

	data, err := base64.StdEncoding.DecodeString(notification.Params.Result.Transaction.Transaction[0])
	if err != nil {
		logger.Log.Printf("[PUMP MONITOR HELIUS]: error b64 decoding %s - %s / %s", txType, sig, err)
		return
	}

	tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(data))
	if err != nil {
		logger.Log.Printf("[PUMP MONITOR HELIUS]: error decoding %s - %s / %s", txType, sig, err)
		return
	}
	preBalances, postBalances := notification.Params.Result.Transaction.Meta.PreBalances, notification.Params.Result.Transaction.Meta.PostBalances
	_, postTokenBalances := notification.Params.Result.Transaction.Meta.PreTokenBalances, notification.Params.Result.Transaction.Meta.PostTokenBalances

	err = pump_monitor_hooks.ParseCreateAndBuy(tx, sig, preBalances, postBalances, postTokenBalances)
	if err != nil {
		logger.Log.Printf("[PUMP MONITOR HELIUS]: error parsing %s - %s / %s", txType, sig, err)
	}
}
