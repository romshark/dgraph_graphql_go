package graphql

import "sync"

// ErrWaitGroup represents an error-wait-group
type ErrWaitGroup interface {
	Inc(delta uint64) uint64
	Dec(delta uint64) uint64
	Fail(err error)
	Wait() error
	IsCompleted() bool
}

type errWaitGroup struct {
	lock      *sync.Mutex
	target    uint64
	err       error
	completed chan struct{}
}

// NewErrWaitGroup creates a new error-wait-group instance
func NewErrWaitGroup(target uint64) ErrWaitGroup {
	return &errWaitGroup{
		lock:      &sync.Mutex{},
		target:    target,
		completed: make(chan struct{}),
	}
}

func (wg *errWaitGroup) IsCompleted() bool {
	select {
	case <-wg.completed:
		return true
	default:
	}
	return false
}

func (wg *errWaitGroup) Inc(delta uint64) uint64 {
	wg.lock.Lock()
	newTarget := wg.target + uint64(delta)
	wg.target = newTarget
	wg.lock.Unlock()
	return newTarget
}

func (wg *errWaitGroup) Dec(delta uint64) uint64 {
	result := uint64(0)
	wg.lock.Lock()
	if wg.target <= delta {
		wg.target = 0
		close(wg.completed)
	} else {
		result = wg.target - delta
		wg.target = result
	}
	wg.lock.Unlock()
	return result
}

func (wg *errWaitGroup) Fail(err error) {
	wg.lock.Lock()
	wg.err = err
	wg.lock.Unlock()
	close(wg.completed)
}

func (wg *errWaitGroup) Wait() error {
	if wg.IsCompleted() {
		wg.lock.Lock()
		err := wg.err
		wg.lock.Unlock()
		return err
	}

	<-wg.completed
	wg.lock.Lock()
	err := wg.err
	wg.lock.Unlock()
	return err
}
