package pump_monitor

import (
	"container/list"
	"sync"
)

var (
	seen     = make(map[string]bool)
	seenMut  sync.Mutex
	queue    = list.New()
	capacity = 100
)

func Seen(sig string) bool {
	seenMut.Lock()
	defer seenMut.Unlock()

	/* check if sig has been seen */
	if _, ok := seen[sig]; ok {
		return true
	}

	/* if sig hasnt been seen add it to the cache */
	seen[sig] = true
	queue.PushBack(sig)

	/* if cache exceeds capacity remove oldest entry */
	if queue.Len() > capacity {
		oldest := queue.Front()
		delete(seen, oldest.Value.(string))
		queue.Remove(oldest)
	}

	return false
}
