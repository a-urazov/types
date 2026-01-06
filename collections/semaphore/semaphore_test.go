package semaphore_test

import (
	"context"
	"testing"
	"time"

	"types/collections/semaphore"
)

func TestSemaphore(t *testing.T) {
	s := semaphore.New(1)

	err := s.Acquire(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan struct{})
	go func() {
		err := s.Acquire(context.Background())
		if err != nil {
			t.Error(err)
		}
		ch <- struct{}{}
	}()

	select {
	case <-ch:
		t.Fatal("should not acquire")
	case <-time.After(100 * time.Millisecond):
	}

	s.Release()

	select {
	case <-ch:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("should acquire")
	}
}
