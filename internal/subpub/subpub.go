package subpub

import (
	"context"
	"sync"
)

type MessageHandler func(msg interface{})

type Subscription interface {
	Unsubscribe()
}

type SubPub interface {
	Subscribe(subject string, cb MessageHandler) (Subscription, error)
	Publish(subject string, msg interface{}) error
	Close(ctx context.Context) error
}

type subPub struct {
	mu          sync.RWMutex
	subjects    map[string]*subject
	closed      bool
	closeSignal chan struct{}
	wg          sync.WaitGroup
}

type subject struct {
	mu          sync.RWMutex
	subscribers map[*subscriber]struct{}
}

type subscriber struct {
	ch      chan interface{}
	cb      MessageHandler
	subject *subject
}

func (s *subscriber) Unsubscribe() {
	s.subject.mu.Lock()
	defer s.subject.mu.Unlock()
	delet(s.subject.subscribers, s)
}

func NewSubPub() SubPub {
	return subPub{mu: sync.RWMutex{}, subjects: map[string]*subject{}, closeSignal: make(chan struct{}), closed: false}
}

func (s *subPub) Subscribe(subject string, cb MessageHandler) (Subscription, error) {}

func (s *subPub) Publish(subject string, msg interface{}) error {}

func (s *subPub) Close(ctx context.Context) error {}
