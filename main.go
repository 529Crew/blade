package main

import (
	"github.com/529Crew/blade/internal/geyser"
	helius_ws "github.com/529Crew/blade/internal/helius/ws"
	pump_monitor "github.com/529Crew/blade/internal/systems/pump/monitor"
)

func main() {
	/* queue up monitors */
	pumpConn := make(chan struct{})
	go pump_monitor.Monitor(pumpConn)

	/* manage geyser connections */
	go geyser.Connect()
	go helius_ws.Connect(pumpConn)

	select {}
}
