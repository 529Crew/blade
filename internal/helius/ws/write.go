package helius_ws

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/FlowGoCrazy/razor/internal/types"
)

var writeMutex = &sync.Mutex{}

func TransactionSubscribe(payload *types.TransactionSubscribePayload) (txSubId int64, e error) {
	if !wsConnected {
		return 0, fmt.Errorf("not connected to geyser")
	}

	/* randomize payload id */
	payload.ID = rand.Intn(10000) + 1

	/* create temp listener to recv tx sub id */
	listenerId, listener := Listen()
	defer Unlisten(listenerId)

	writeMutex.Lock()
	defer writeMutex.Unlock()

	err := wsConnection.WriteJSON(payload)
	if err != nil {
		return 0, fmt.Errorf("error writing transactionSubscribe to geyser: %v", err)
	}

	/* wait for tx sub id */
	for {
		select {
		case msg := <-listener:
			/* attempt to decode to tx sub response */
			var txSubResp types.TransactionSubscribeResponse
			err = json.Unmarshal(msg, &txSubResp)
			if err != nil {
				continue
			}

			/* ensure response is a tx sub response and matches tx sub id */
			if txSubResp.ID != payload.ID || txSubResp.Result == 0 {
				continue
			}

			return txSubResp.Result, nil
		case <-time.After(5 * time.Second):
			return 0, fmt.Errorf("timeout waiting for transactionSubscribe response")
		}
	}
}
