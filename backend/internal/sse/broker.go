package sse

import (
	"sync"

	"github.com/secscan/secscan/backend/internal/domain"
)

type Broker struct {
	mu          sync.RWMutex
	subscribers map[string]map[chan domain.ProgressEvent]struct{}
}

func NewBroker() *Broker {
	return &Broker{subscribers: make(map[string]map[chan domain.ProgressEvent]struct{})}
}

func (b *Broker) Subscribe(scanID string) (<-chan domain.ProgressEvent, func()) {
	ch := make(chan domain.ProgressEvent, 32)

	b.mu.Lock()
	if b.subscribers[scanID] == nil {
		b.subscribers[scanID] = make(map[chan domain.ProgressEvent]struct{})
	}
	b.subscribers[scanID][ch] = struct{}{}
	b.mu.Unlock()

	unsubscribe := func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		if subs, ok := b.subscribers[scanID]; ok {
			if _, exists := subs[ch]; exists {
				delete(subs, ch)
				close(ch)
			}
			if len(subs) == 0 {
				delete(b.subscribers, scanID)
			}
		}
	}

	return ch, unsubscribe
}

func (b *Broker) Publish(event domain.ProgressEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for ch := range b.subscribers[event.ScanID] {
		select {
		case ch <- event:
		default:
		}
	}
}
