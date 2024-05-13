package helius_ws

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/FlowGoCrazy/razor/internal/logger"
	"github.com/gorilla/websocket"
)

const retryInterval = 5 * time.Second

var wsConnection *websocket.Conn
var wsConnected bool

func Connect(connChans ...chan struct{}) error {
	var err error

	for {
		u, parseErr := url.Parse(fmt.Sprintf("wss://atlas-mainnet.helius-rpc.com?api-key=%s", os.Getenv("HELIUS_API_KEY")))
		if parseErr != nil {
			return fmt.Errorf("error parsing ws url: %v", parseErr)
		}

		logger.Log.Println("[HELIUS WS] connecting")
		wsConnection, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			logger.Log.Println("[HELIUS WS] failed to connect, retrying...")
			time.Sleep(retryInterval)
			continue
		}
		logger.Log.Println("[HELIUS WS] connected")
		wsConnected = true

		/* alert monitors that new connection was established */
		for _, conn := range connChans {
			go func(c chan struct{}) {
				c <- struct{}{}
			}(conn)
		}

		for {
			_, msg, err := wsConnection.ReadMessage()
			if err != nil {
				logger.Log.Printf("[HELIUS WS] disconnected: %s\n", err)
				wsConnected = false
				wsConnection.Close()
				break
			}

			/* emit message to all listeners */
			listenersMutex.Lock()
			for _, listener := range wsListeners {
				go func(l chan []byte) {
					l <- msg
				}(listener)
			}
			listenersMutex.Unlock()
		}
	}
}
