package signal

import (
	"testing"
	"time"
)

// TestSignal tests the basic functionality of the Signal type
func TestSignal(t *testing.T) {
	sig := New[int]()

	// Subscribe to the signal
	ch, unsubscribe := sig.Subscribe(10)

	// Test emitting a value
	sig.Emit(42)

	// Check if we receive the value
	select {
	case val := <-ch:
		if val != 42 {
			t.Errorf("Expected 42, got %d", val)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Did not receive value within timeout")
	}

	// Test unsubscribe
	unsubscribe()

	// Try to send another value after unsubscribe
	sig.Emit(100)

	// Channel should be closed now, so we should not receive anything
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("Channel should be closed after unsubscribe")
		}
	default:
		// Channel is closed, which is expected
	}
}

// TestSignalConcurrent tests the Signal type with concurrent operations
func TestSignalConcurrent(t *testing.T) {
	sig := New[string]()

	// Create multiple subscribers
	ch1, unsub1 := sig.Subscribe(10)
	ch2, unsub2 := sig.Subscribe(10)
	ch3, unsub3 := sig.Subscribe(10)

	// Emit value concurrently
	go func() {
		sig.Emit("test")
	}()

	// Receive on all channels
	received := make([]string, 3)
	done := make(chan bool, 3)

	go func() {
		received[0] = <-ch1
		done <- true
	}()

	go func() {
		received[1] = <-ch2
		done <- true
	}()

	go func() {
		received[2] = <-ch3
		done <- true
	}()

	// Wait for all receives
	for i := 0; i < 3; i++ {
		select {
		case <-done:
		case <-time.After(100 * time.Millisecond):
			t.Error("Did not receive value within timeout")
		}
	}

	// Verify all received values
	for i := 0; i < 3; i++ {
		if received[i] != "test" {
			t.Errorf("Expected 'test', got '%s' at index %d", received[i], i)
		}
	}

	// Unsubscribe all
	unsub1()
	unsub2()
	unsub3()
}

// TestBufferedSignal tests the BufferedSignal type
func TestBufferedSignal(t *testing.T) {
	bufSig := NewBuffered[int](5)

	// Subscribe to the buffered signal
	ch, unsubscribe := bufSig.Subscribe(10)

	// Emit multiple values
	for i := 1; i <= 3; i++ {
		bufSig.Emit(i)
	}

	// Check if we receive all values
	for i := 1; i <= 3; i++ {
		select {
		case val := <-ch:
			if val != i {
				t.Errorf("Expected %d, got %d", i, val)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Did not receive value %d within timeout", i)
		}
	}

	// Check buffer contents
	buffer := bufSig.GetBuffer()
	if len(buffer) != 3 {
		t.Errorf("Expected buffer length 3, got %d", len(buffer))
	}

	for i, val := range buffer {
		if val != i+1 {
			t.Errorf("Expected %d at index %d, got %d", i+1, i, val)
		}
	}

	unsubscribe()
}

// TestBufferedSignalBufferLimit tests that the buffer limit is respected
func TestBufferedSignalBufferLimit(t *testing.T) {
	bufSig := NewBuffered[int](3) // Buffer limit of 3

	// Add more values than the buffer limit
	for i := 1; i <= 5; i++ {
		bufSig.Emit(i)
	}

	// Check buffer contents - should only have the last 3 values (3, 4, 5)
	buffer := bufSig.GetBuffer()
	if len(buffer) != 3 {
		t.Errorf("Expected buffer length 3, got %d", len(buffer))
	}

	expected := []int{3, 4, 5}
	for i, val := range buffer {
		if val != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, val)
		}
	}
}

// TestBufferedSignalClearBuffer tests clearing the buffer
func TestBufferedSignalClearBuffer(t *testing.T) {
	bufSig := NewBuffered[int](5)

	// Add some values
	for i := 1; i <= 3; i++ {
		bufSig.Emit(i)
	}

	// Verify buffer has values
	buffer := bufSig.GetBuffer()
	if len(buffer) != 3 {
		t.Errorf("Expected buffer length 3 before clear, got %d", len(buffer))
	}

	// Clear the buffer
	bufSig.ClearBuffer()

	// Verify buffer is empty
	buffer = bufSig.GetBuffer()
	if len(buffer) != 0 {
		t.Errorf("Expected buffer length 0 after clear, got %d", len(buffer))
	}
}

// TestSignalEmitSync tests the EmitSync method
func TestSignalEmitSync(t *testing.T) {
	sig := New[int]()

	// Subscribe to the signal
	ch, unsubscribe := sig.Subscribe(10)
	defer unsubscribe()

	// Use EmitSync to send value
	sig.EmitSync(123)

	// Check if we receive the value
	select {
	case val := <-ch:
		if val != 123 {
			t.Errorf("Expected 123, got %d", val)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Did not receive value within timeout")
	}
}

// TestSignalBroadcast tests the Broadcast method
func TestSignalBroadcast(t *testing.T) {
	sig := New[string]()

	// Subscribe multiple listeners
	ch1, unsub1 := sig.Subscribe(10)
	ch2, unsub2 := sig.Subscribe(10)
	ch3, unsub3 := sig.Subscribe(10)
	defer unsub1()
	defer unsub2()
	defer unsub3()

	// Use Broadcast to send value
	sig.Broadcast("broadcast_test")

	// Check if all channels receive the value
	received := make([]string, 3)
	done := make(chan bool, 3)

	go func() {
		received[0] = <-ch1
		done <- true
	}()

	go func() {
		received[1] = <-ch2
		done <- true
	}()

	go func() {
		received[2] = <-ch3
		done <- true
	}()

	// Wait for all receives
	for i := 0; i < 3; i++ {
		select {
		case <-done:
		case <-time.After(100 * time.Millisecond):
			t.Error("Did not receive value within timeout")
		}
	}

	// Verify all received values
	for i := 0; i < 3; i++ {
		if received[i] != "broadcast_test" {
			t.Errorf("Expected 'broadcast_test', got '%s' at index %d", received[i], i)
		}
	}
}
