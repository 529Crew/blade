package main

import (
	"github.com/529Crew/blade/internal/geyser"
	helius_ws "github.com/529Crew/blade/internal/helius/ws"
	pump_monitor "github.com/529Crew/blade/internal/systems/pump/monitor"
)

func main() {
	/* queue up helius monitors */
	pumpConnHelius := make(chan struct{})
	go pump_monitor.Monitor(pumpConnHelius)

	/* manage helius connection */
	go helius_ws.Connect(pumpConnHelius)

	/* queue up geyser monitors */
	pumpConnGeyser := make(chan struct{})
	go pump_monitor.MonitorGeyser(pumpConnGeyser)

	/* manage geyser connection */
	go geyser.Connect(pumpConnGeyser)

	select {}
}
