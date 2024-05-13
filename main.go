package main

import (
	helius_ws "github.com/529Crew/blade/internal/helius/ws"
)

func main() {
	/* manage geyser connection */
	go helius_ws.Connect()

	select {}
}
