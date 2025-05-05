package subpub

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	bus := NewSubPub()
	received := make(chan bool)
	sub, _ := bus.Subscribe("news", func(msg interface{}) {
		if msg != "hello" {
			t.Error("Unexpected message")
		}
		received <- true
	})
	if sub == nil {
		t.Error("Subscriber is nil")
	}
	defer sub.Unsubscribe()
	if err := bus.Publish("news", "hello"); err != nil {
		t.Fatal(err)
	}
	select {
	case <-received:
		// OK
	case <-time.After(1 * time.Second):
		t.Fatal("Message not received")
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
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub, _ := bus.Subscribe("news", func(msg interface{}) {})
			defer sub.Unsubscribe()
		}()
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
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
