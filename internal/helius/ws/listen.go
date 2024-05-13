package helius_ws

import (
	"sync"

	"github.com/google/uuid"
)

var wsListeners = make(map[string]chan []byte)
var listenersMutex = sync.Mutex{}

func Listen() (string, chan []byte) {
	listenersMutex.Lock()
	defer listenersMutex.Unlock()

	listenerId := uuid.New().String()

	listener := make(chan []byte)
	wsListeners[listenerId] = listener

	return listenerId, listener
}

func Unlisten(listenerId string) {
	listenersMutex.Lock()
	defer listenersMutex.Unlock()

	delete(wsListeners, listenerId)
}
