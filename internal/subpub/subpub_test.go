package subpub

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	tests := []struct {
		name        string
		subject     string
		message     interface{}
		expectError bool
	}{
		{"Good case", "news", "hello", false},
		{"Empty subject", "", "test", false},
		{"Not string message", "news", 1, false},
	}
	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := NewSubPub()
			received := make(chan bool)
			sub, err := bus.Subscribe(tt.subject, func(msg interface{}) {
				if msg != tt.message {
					t.Errorf("Expected message %v, got %v", tt.message, msg)
				}
				received <- true
			})
			if err != nil {
				t.Fatalf("Subscribe error: %v", err)
			}
			if sub == nil {
				t.Error("Subscriber is nil")
			}
			defer sub.Unsubscribe()

			err = bus.Publish(tt.subject, tt.message)
			if (err != nil) != tt.expectError {
				t.Errorf("Publish error: %v, want error: %v", err, tt.expectError)
			}

			select {
			case <-received:
				// OK
			case <-time.After(1 * time.Second):
				if !tt.expectError {
					t.Error("Message not received")
				}
			}
		})
	}
}

func TestUnsubscribe(t *testing.T) {
	bus := NewSubPub()
	received := false
	sub, _ := bus.Subscribe("news", func(msg interface{}) {
		received = true
	})
	sub.Unsubscribe()
	bus.Publish("news", "test")
	if received {
		t.Error("Received message after unsubscribe")
	}
}

func TestClose(t *testing.T) {
	bus := NewSubPub()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	bus.Subscribe("news", func(msg interface{}) {})
	if err := bus.Close(ctx); err != nil {
		t.Fatal(err)
	}
	if err := bus.Publish("news", "test"); err == nil {
		t.Error("Expected error after Close")
	}
	if _, err := bus.Subscribe("post", func(msg interface{}) {}); err == nil {
		t.Error("Expected error after Close")
	}
}

func TestConcurrency(t *testing.T) {
	bus := NewSubPub()
	var wg sync.WaitGroup
	const numOps = 100
	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub, _ := bus.Subscribe("news", func(msg interface{}) {})
			defer sub.Unsubscribe()
			bus.Publish("news", "test")
		}()
	}
	wg.Wait()
}

func TestPublishToNonExistentSubject(t *testing.T) {
	bus := NewSubPub()
	err := bus.Publish("nonexistent", "test")
	if err != nil {
		t.Errorf("Expected no error when publishing to non-existent subject, got %v", err)
	}
}

func TestCloseAlreadyClosed(t *testing.T) {
	bus := NewSubPub()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if err := bus.Close(ctx); err != nil {
		t.Fatalf("First Close error: %v", err)
	}
	if err := bus.Close(ctx); err != nil {
		t.Errorf("Second Close error: %v, expected nil", err)
	}
}

func TestCloseWithActiveSubscribers(t *testing.T) {
	bus := NewSubPub()
	block := make(chan struct{})
	received := make(chan struct{})
	bus.Subscribe("block", func(msg interface{}) {
		<-block
		received <- struct{}{}
	})
	bus.Publish("block", "message")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	err := bus.Close(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}
	close(block)
	<-received
}
