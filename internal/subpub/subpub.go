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
	delete(s.subject.subscribers, s)
	close(s.ch)
}

func NewSubPub() SubPub {
	return &subPub{
		subjects:    make(map[string]*subject),
		closeSignal: make(chan struct{}),
	}
}

func (s *subPub) Subscribe(subjectName string, cb MessageHandler) (Subscription, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	subj, ok := s.subjects[subjectName]
	if !ok {
		subj = &subject{
			subscribers: make(map[*subscriber]struct{}),
		}
		s.subjects[subjectName] = subj
	}
	sub := &subscriber{
		ch:      make(chan interface{}, 10),
		cb:      cb,
		subject: subj,
	}
	subj.mu.Lock()
	subj.subscribers[sub] = struct{}{}
	subj.mu.Unlock()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case msg, ok := <-sub.ch:
				if !ok {
					return
				}
				sub.cb(msg)
			case <-s.closeSignal:
				return
			}
		}
	}()
	return sub, nil
}

func (s *subPub) Publish(subject string, msg interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	subj, ok := s.subjects[subject]
	if !ok {
		return nil
	}
	subj.mu.RLock()
	defer subj.mu.RUnlock()
	for sub := range subj.subscribers {
		sub.ch <- msg
	}
	return nil
}

func (s *subPub) Close(ctx context.Context) error {}
