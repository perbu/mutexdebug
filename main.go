package mutexdebug

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type DebugMutex struct {
	mux      sync.Mutex
	done     chan bool
	timeout  time.Duration
	Warnings atomic.Uint64
	warn     bool
}

func NewDebugMutex(timeout time.Duration, warn bool) *DebugMutex {
	return &DebugMutex{
		mux:     sync.Mutex{},
		done:    make(chan bool),
		timeout: timeout,
		warn:    warn,
	}
}

func (dm *DebugMutex) Lock() {
	// identify the caller:
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("Could not identify caller")
	}
	// fmt.Printf("Locking mutex from %s:%d\n", file, line)

	dm.mux.Lock()
	go func() {
		timer := time.NewTimer(dm.timeout)
		for {
			select {
			case <-dm.done:
				return
			case <-timer.C:
				// print a warning if configured to do so
				dm.Warnings.Add(1)
				timer.Stop()
				if !dm.warn {
					continue
				}
				fmt.Printf("Mutex locked for too long from %s:%d\n", file, line)
			}
		}
	}()
}

func (dm *DebugMutex) Unlock() {
	dm.mux.Unlock()
	dm.done <- true

}
