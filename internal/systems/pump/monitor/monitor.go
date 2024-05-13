package pump_monitor

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/529Crew/blade/internal/logger"
	"github.com/529Crew/blade/internal/types"

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
	// err = parseCreate(&notification)
	case strings.Contains(combinedLogs, "Instruction: Create") && !strings.Contains(combinedLogs, "Instruction: Buy"):
		txType = "create"
	// 	err = parseCreate(&notification)
	case !strings.Contains(combinedLogs, "Instruction: Create") && strings.Contains(combinedLogs, "Instruction: Buy"):
		txType = "buy"
		// err = parseBuy(&notification)
	}
	// if err != nil {
	// 	logger.Log.Printf("[PUMP MONITOR HELIUS]: %s - %s / %s", txType, sig, err)
	// }

	if txType != "" && txType != "buy" {
		logger.Log.Printf("[PUMP MONITOR HELIUS]: %s / %s", txType, sig)
	}
}
