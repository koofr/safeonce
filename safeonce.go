package safeonce

import (
	"sync"
	"sync/atomic"
)

// SafeOnce is an object that will perform one action (until it succeeds).
type SafeOnce struct {
	m    sync.Mutex
	done uint32
}

// Do calls the function f if and only if Do is being called for the
// first time for this instance of Once or if previous calls returned error.
// In other words, given
// 	var onceSafe OnceSafe
// if once.Do(f) is called multiple times, only the first call that returns
// without error will invoke f, even if f has a different value in each
// invocation.  A new instance of SafeOnce is required for each function
// to execute.
//
// Because no call to Do returns until the one call to f returns, if f causes
// Do to be called, it will deadlock.
//
// If f panics, Do considers it to have returned; future calls of Do return
// without calling f.
//
func (o *SafeOnce) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) == 1 {
		return nil
	}
	// Slow-path.
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		var err error
		defer func() {
			if err == nil {
				atomic.StoreUint32(&o.done, 1)
			}
		}()
		err = f()
		return err
	}
	return nil
}
