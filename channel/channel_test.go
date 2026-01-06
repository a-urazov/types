// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package channel_test

import (
	"context"
	"testing"

	"types/channel"
)

func TestChannelSendAndReceive(t *testing.T) {
	ch := channel.New[int](0)
	go func() {
		err := ch.Send(context.Background(), 1)
		if err != nil {
			t.Errorf("Send error: %v", err)
		}
	}()

	val, err := ch.Receive(context.Background())
	if err != nil {
		t.Fatalf("Receive error: %v", err)
	}
	if val != 1 {
		t.Errorf("Expected 1, got %d", val)
	}
}

func TestChannelBuffered(t *testing.T) {
	ch := channel.New[int](1)
	err := ch.Send(context.Background(), 1)
	if err != nil {
		t.Fatalf("Send error: %v", err)
	}

	val, err := ch.Receive(context.Background())
	if err != nil {
		t.Fatalf("Receive error: %v", err)
	}
	if val != 1 {
		t.Errorf("Expected 1, got %d", val)
	}
}

func TestChannelClose(t *testing.T) {
	ch := channel.New[int](0)
	ch.Close()

	err := ch.Send(context.Background(), 1)
	if err != channel.ErrClosedChannel {
		t.Errorf("Expected ErrClosedChannel, got %v", err)
	}

	_, err = ch.Receive(context.Background())
	if err != channel.ErrClosedChannel {
		t.Errorf("Expected ErrClosedChannel, got %v", err)
	}
}

func TestChannelContextCancellation(t *testing.T) {
	ch := channel.New[int](0)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := ch.Send(ctx, 1)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}

	_, err = ch.Receive(ctx)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestChannelRange(t *testing.T) {
	ch := channel.New[int](3)
	ch.Send(context.Background(), 1)
	ch.Send(context.Background(), 2)
	ch.Send(context.Background(), 3)
	ch.Close()

	var sum int
	ch.Range(func(val int) bool {
		sum += val
		return true
	})

	if sum != 6 {
		t.Errorf("Expected sum 6, got %d", sum)
	}
}

func BenchmarkChannel(b *testing.B) {
	b.Run("Unbuffered", func(b *testing.B) {
		ch := channel.New[int](0)
		go func() {
			for i := 0; i < b.N; i++ {
				ch.Send(context.Background(), i)
			}
			ch.Close()
		}()
		// Consume all values from the channel (benchmark only needs to measure send performance)
		for range ch.Unwrap() {
			// Intentionally empty - just consuming values
		}
	})

	b.Run("Buffered", func(b *testing.B) {
		ch := channel.New[int](128)
		go func() {
			for i := 0; i < b.N; i++ {
				ch.Send(context.Background(), i)
			}
			ch.Close()
		}()
		// Consume all values from the channel (benchmark only needs to measure send performance)
		for range ch.Unwrap() {
			// Intentionally empty - just consuming values
		}
	})
}
