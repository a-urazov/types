// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package channel_test

import (
	"context"
	"testing"
	"time"

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

func TestChannelSelect(t *testing.T) {
	// Test basic receive case
	ch1 := channel.New[int]()
	go func() {
		time.Sleep(10 * time.Millisecond)
		ch1.Send(context.Background(), 42)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	case1 := &channel.Case[int]{
		Chan:  ch1,
		Send:  false, // receive operation
		Value: 0,
	}

	index, val, err := ch1.Select(ctx, case1)
	if err != nil {
		t.Fatalf("Select error: %v", err)
	}
	if index != 0 {
		t.Errorf("Expected case index 0, got %d", index)
	}
	if val != 42 {
		t.Errorf("Expected value 42, got %d", val)
	}

	// Test basic send case
	ch2 := channel.New[int]()
	case2 := &channel.Case[int]{
		Chan:  ch2,
		Send:  true, // send operation
		Value: 100,
	}

	go func() {
		received, err := ch2.Receive(context.Background())
		if err != nil {
			t.Errorf("Receive error: %v", err)
			return
		}
		if received != 100 {
			t.Errorf("Expected received value 100, got %d", received)
		}
	}()

	index, val, err = ch2.Select(ctx, case2)
	if err != nil {
		t.Fatalf("Select error: %v", err)
	}
	if index != 0 {
		t.Errorf("Expected case index 0, got %d", index)
	}
	if val != 100 {
		t.Errorf("Expected value 100, got %d", val)
	}

	// Test multiple cases with context cancellation
	ch3 := channel.New[int]()
	ch4 := channel.New[int]() // Use same type

	case3 := &channel.Case[int]{
		Chan:  ch3,
		Send:  false,
		Value: 0,
	}
	case4 := &channel.Case[int]{
		Chan:  ch4,
		Send:  false,
		Value: 0,
	}

	shortCtx, cancel2 := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel2()

	index, _, err = ch3.Select(shortCtx, case3, case4)
	if err == nil {
		t.Error("Expected context cancellation error")
	}
	if index != -1 {
		t.Errorf("Expected index -1 on context cancellation, got %d", index)
	}
}

func TestChannelTryOperations(t *testing.T) {
	// Test TrySend on buffered channel
	ch := channel.New[int](2)

	// Should be able to send to buffered channel
	ok := ch.TrySend(10)
	if !ok {
		t.Error("Expected TrySend to succeed on buffered channel")
	}

	ok = ch.TrySend(20)
	if !ok {
		t.Error("Expected TrySend to succeed on buffered channel")
	}

	// Buffer is full, should fail
	ok = ch.TrySend(30)
	if ok {
		t.Error("Expected TrySend to fail when buffer is full")
	}

	// Test TryReceive
	val, ok := ch.TryReceive()
	if !ok || val != 10 {
		t.Errorf("Expected TryReceive to return (10, true), got (%d, %v)", val, ok)
	}

	// Test TryReceive on empty channel
	ch2 := channel.New[int]()
	_, ok = ch2.TryReceive()
	if ok {
		t.Error("Expected TryReceive to fail on empty channel")
	}
}

func TestChannelDrain(t *testing.T) {
	ch := channel.New[int](5)

	// Fill the channel
	for i := 0; i < 5; i++ {
		ch.Send(context.Background(), i*10)
	}

	// Drain it
	count := ch.Drain()
	if count != 5 {
		t.Errorf("Expected to drain 5 values, got %d", count)
	}

	// Channel should still be open, just empty
	// Try to receive with timeout to verify it's empty
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := ch.Receive(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded for empty channel, got %v", err)
	}

	// Send to verify channel is still functional
	err = ch.Send(context.Background(), 100)
	if err != nil {
		t.Errorf("Expected Send to succeed on drained channel, got %v", err)
	}

	// Receive the value we just sent
	val, err := ch.Receive(context.Background())
	if err != nil || val != 100 {
		t.Errorf("Expected to receive 100, got %d with error %v", val, err)
	}
}

func TestChannelMerge(t *testing.T) {
	ch1 := channel.New[int]()
	ch2 := channel.New[int]()
	ch3 := channel.New[int]()

	// Send values to each channel
	go func() {
		ch1.Send(context.Background(), 1)
		ch1.Send(context.Background(), 2)
		ch1.Close()
	}()

	go func() {
		ch2.Send(context.Background(), 3)
		ch2.Send(context.Background(), 4)
		ch2.Close()
	}()

	go func() {
		ch3.Send(context.Background(), 5)
		ch3.Send(context.Background(), 6)
		ch3.Close()
	}()

	// Merge the channels
	merged := channel.Merge(ch1, ch2, ch3)

	// Collect all values
	values := []int{}
	for {
		val, err := merged.Receive(context.Background())
		if err != nil {
			break // channel closed
		}
		values = append(values, val)
	}

	if len(values) != 6 {
		t.Errorf("Expected 6 values from merged channels, got %d", len(values))
	}
}

func TestChannelFanOut(t *testing.T) {
	input := channel.New[int]()
	output1 := channel.New[int]()
	output2 := channel.New[int]()

	// Start fan out
	channel.FanOut(input, output1, output2)

	// Send value to input
	input.Send(context.Background(), 42)
	input.Close()

	// Receive from both outputs
	val1, err1 := output1.Receive(context.Background())
	if err1 != nil {
		t.Errorf("Error receiving from output1: %v", err1)
	}

	val2, err2 := output2.Receive(context.Background())
	if err2 != nil {
		t.Errorf("Error receiving from output2: %v", err2)
	}

	if val1 != 42 || val2 != 42 {
		t.Errorf("Expected both outputs to receive 42, got %d and %d", val1, val2)
	}
}

func TestChannelBatchOperations(t *testing.T) {
	// Test BatchSend
	ch := channel.New[int](10)
	valuesToSend := []int{1, 2, 3, 4, 5}

	sentCount := ch.BatchSend(context.Background(), valuesToSend)
	if sentCount != len(valuesToSend) {
		t.Errorf("Expected to send %d values, sent %d", len(valuesToSend), sentCount)
	}

	// Test BatchReceive
	receivedValues, err := ch.BatchReceive(context.Background(), 3)
	if err != nil {
		t.Errorf("BatchReceive error: %v", err)
	}
	if len(receivedValues) != 3 {
		t.Errorf("Expected to receive 3 values, got %d", len(receivedValues))
	}
	for i, val := range receivedValues {
		if val != i+1 {
			t.Errorf("Expected value %d at position %d, got %d", i+1, i, val)
		}
	}
}

func TestChannelPipelineOperations(t *testing.T) {
	// Test Take
	source := channel.New[int]()
	go func() {
		for i := 1; i <= 5; i++ {
			source.Send(context.Background(), i)
		}
		source.Close()
	}()

	take3 := channel.Take(source, 3)
	takenValues := []int{}
	for {
		val, err := take3.Receive(context.Background())
		if err != nil {
			break
		}
		takenValues = append(takenValues, val)
	}

	if len(takenValues) != 3 {
		t.Errorf("Expected to take 3 values, got %d", len(takenValues))
	}

	// Test Skip
	source2 := channel.New[int](10) // Используем буферизованный канал
	go func() {
		for i := 1; i <= 5; i++ {
			source2.Send(context.Background(), i)
		}
		source2.Close()
	}()

	skip2 := channel.Skip(source2, 2)
	skippedValues := []int{}
	for {
		val, err := skip2.Receive(context.Background())
		if err != nil {
			break
		}
		skippedValues = append(skippedValues, val)
	}

	expectedSkipped := []int{3, 4, 5}
	if len(skippedValues) != len(expectedSkipped) {
		t.Errorf("Expected to skip to %d values, got %d", len(expectedSkipped), len(skippedValues))
	}
}

func TestChannelFilterMap(t *testing.T) {
	// Test Filter
	source := channel.New[int]()
	for i := 1; i <= 5; i++ {
		source.Send(context.Background(), i)
	}
	source.Close()

	evenFilter := channel.Filter(source, func(x int) bool { return x%2 == 0 })
	evenValues := []int{}
	for {
		val, err := evenFilter.Receive(context.Background())
		if err != nil {
			break
		}
		evenValues = append(evenValues, val)
	}

	if len(evenValues) != 2 { // Should be [2, 4]
		t.Errorf("Expected 2 even values, got %d", len(evenValues))
	}

	// Test Map
	source2 := channel.New[int]()
	for i := 1; i <= 3; i++ {
		source2.Send(context.Background(), i)
	}
	source2.Close()

	squaredMap := channel.Map(source2, func(x int) int { return x * x })
	squaredValues := []int{}
	for {
		val, err := squaredMap.Receive(context.Background())
		if err != nil {
			break
		}
		squaredValues = append(squaredValues, val)
	}

	if len(squaredValues) != 3 { // Should be [1, 4, 9]
		t.Errorf("Expected 3 squared values, got %d", len(squaredValues))
		ch := channel.New[int]()

		// Test ConditionalSend with condition that passes
		err := ch.ConditionalSend(context.Background(), 42, func(x int) bool { return x > 40 })
		if err != nil {
			t.Errorf("Expected ConditionalSend to succeed with passing condition, got error: %v", err)
		}

		// Receive the value to clear the channel
		val, err := ch.Receive(context.Background())
		if err != nil || val != 42 {
			t.Errorf("Expected to receive 42, got %d with error %v", val, err)
		}

		// Test ConditionalSend with condition that fails
		err = ch.ConditionalSend(context.Background(), 10, func(x int) bool { return x > 40 })
		if err != nil {
			t.Errorf("Expected ConditionalSend to return nil when condition fails, got error: %v", err)
		}

		// Channel should remain empty
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()
		_, err = ch.Receive(ctx)
		if err != context.DeadlineExceeded {
			t.Errorf("Expected timeout error for empty channel, got %v", err)
		}

		// Test ConditionalReceive
		ch.Send(context.Background(), 5)
		ch.Send(context.Background(), 15)
		ch.Send(context.Background(), 25)

		// Receive only values greater than 10
		condition := func(x int) bool { return x > 10 }
		val, err = ch.ConditionalReceive(context.Background(), condition)
		if err != nil || val != 15 {
			t.Errorf("Expected to receive 15 (first value > 10), got %d with error %v", val, err)
		}
	}

}
