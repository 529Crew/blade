package pump_monitor

import (
	"container/list"
	"sync"
)

var (
	queue    = list.New()
	capacity = 100
	mutex    sync.Mutex
)

func Seen(sig string) bool {
	mutex.Lock()
	defer mutex.Unlock()

	/* check if the sig has already been seen */
	for e := queue.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == sig {
			return true
		}
	}

	/* if it hasnt been seen add it to the cache */
	queue.PushBack(sig)

	/* if cache exceeds capacity remove the oldest entry */
	if queue.Len() > capacity {
		queue.Remove(queue.Front())
	}

	return false
}
