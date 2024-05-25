package main

import (
	"log"
	"os"
	"time"

	"github.com/529Crew/blade/internal/geyser"
	helius_ws "github.com/529Crew/blade/internal/helius/ws"
	pump_monitor "github.com/529Crew/blade/internal/systems/pump/monitor"
	self_monitor "github.com/529Crew/blade/internal/systems/self/monitor"
)

var expiry int64 = 1717128000

func main() {
	/* check expiry */
	go func() {
		for {
			if expiry < time.Now().Unix() {
				log.Fatal("binary expired, contact me! exiting in 5 seconds...")
				time.Sleep(5 * time.Second)
				os.Exit(1)
			}

			time.Sleep(30 * time.Second)
		}
	}()

	/* queue up helius monitors */
	pumpConnHelius := make(chan struct{})
	go pump_monitor.Monitor(pumpConnHelius)

	selfConnHelius := make(chan struct{})
	go self_monitor.Monitor(selfConnHelius)

	/* manage helius connection */
	go helius_ws.Connect(pumpConnHelius, selfConnHelius)

	/* queue up geyser monitors */
	pumpConnGeyser := make(chan struct{})
	go pump_monitor.MonitorGeyser(pumpConnGeyser)

	/* manage geyser connection */
	go geyser.Connect(pumpConnGeyser)

	select {}
}
