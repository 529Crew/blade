package geyser

import (
	"sync"

	"github.com/google/uuid"

	pb "github.com/529Crew/blade/internal/geyser/proto"
)

/* map [ label ] map [ listenerId ] chan */
var labelListeners = make(map[string]map[string]chan *pb.SubscribeUpdate)
var labelListenersMutex = &sync.Mutex{}

func AlertListeners(subUpdate *pb.SubscribeUpdate) {
	if len(subUpdate.Filters) < 1 {
		return
	}
	label := subUpdate.Filters[0]

	/* broadcast to all listeners subscribed to label */
	labelListenersMutex.Lock()
	for _, listener := range labelListeners[label] {
		go func(l chan *pb.SubscribeUpdate) {
			l <- subUpdate
		}(listener)
	}
	labelListenersMutex.Unlock()
}

func Listen(label string) (string, chan *pb.SubscribeUpdate) {
	labelListenersMutex.Lock()
	defer labelListenersMutex.Unlock()

	listenerId := uuid.New().String()

	listener := make(chan *pb.SubscribeUpdate)

	if _, ok := labelListeners[label]; !ok {
		labelListeners[label] = make(map[string]chan *pb.SubscribeUpdate)
	}
	labelListeners[label][listenerId] = listener

	return listenerId, listener
}

func Unlisten(label string, listenerId string) {
	labelListenersMutex.Lock()
	defer labelListenersMutex.Unlock()

	delete(labelListeners[label], listenerId)
}
