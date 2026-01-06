package semaphore

import "context"

// Semaphore is a semaphore implementation.
type Semaphore struct {
	ch chan struct{}
}

// New creates a new Semaphore with the given capacity.
func New(capacity int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, capacity),
	}
}

// Acquire acquires a semaphore.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.ch <- struct{}{}:
		return nil
	}
}

// Release releases a semaphore.
func (s *Semaphore) Release() {
	<-s.ch
}
