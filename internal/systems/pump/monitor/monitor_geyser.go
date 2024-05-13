package pump_monitor

import (
	"strings"

	"github.com/529Crew/blade/internal/geyser"
	"github.com/529Crew/blade/internal/logger"
	"github.com/btcsuite/btcd/btcutil/base58"

	pb "github.com/529Crew/blade/internal/geyser/proto"
)

func MonitorGeyser(connected chan struct{}) {
	var listenerId string
	var listener chan *pb.SubscribeUpdate

	for {
		select {
		case <-connected:
			/* handle possible existing subscription */
			if listenerId != "" {
				geyser.UnsubscribeTransactions("pump", listenerId)
			}

			/* subscribe to pump transactions */
			lid, l, err := geyser.SubscribeTransactions(
				"pump",
				&pb.SubscribeRequestFilterTransactions{
					Vote:   new(bool),
					Failed: new(bool),
					AccountRequired: []string{
						"6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P",
					},
				},
			)
			if err != nil {
				logger.Log.Printf("[PUMP MONITOR GEYSER] subscription error: %s", err)
			}

			listenerId = lid
			listener = l

			logger.Log.Println("[PUMP MONITOR GEYSER] activated")
		case msg := <-listener:
			go sortGeyser(msg.GetTransaction())
		}
	}
}

func sortGeyser(msg *pb.SubscribeUpdateTransaction) {
	combinedLogs := strings.Join(msg.Transaction.Meta.LogMessages, " ")
	sig := base58.Encode(msg.Transaction.Signature)

	var txType string
	// var err error

	switch {
	case strings.Contains(combinedLogs, "Instruction: Create") && strings.Contains(combinedLogs, "Instruction: Buy"):
		txType = "create + buy"
	// err = parseCreate(&notification)
	case strings.Contains(combinedLogs, "Instruction: Create") && !strings.Contains(combinedLogs, "Instruction: Buy"):
		txType = "create"
	// 	err = parseCreate(&notification)
	case !strings.Contains(combinedLogs, "Instruction: Create") && strings.Contains(combinedLogs, "Instruction: Buy"):
		txType = "buy"
		// err = parseBuy(msg)
	}
	// if err != nil {
	// 	logger.Log.Printf("[PUMP MONITOR GEYSER] %s - %s: %s", txType, sig[:5], err)
	// }

	if txType != "" && txType != "buy" {
		logger.Log.Printf("[PUMP MONITOR GEYSER] %s: %s\n", txType, sig)
	}
}
