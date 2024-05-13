package pump_monitor

import "sync"

var seen = make(map[string]bool)
var seenMut = sync.Mutex{}

func Seen(sig string) bool {
	seenMut.Lock()
	defer seenMut.Unlock()

	/* if we have seen this sig */
	_, ok := seen[sig]
	if ok {
		return true
	}

	/* if we havent seen this sig */
	seen[sig] = true
	return false
}
