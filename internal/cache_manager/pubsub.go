package cache_manager

import "sync"

type PubSub struct {
	mu          sync.RWMutex
	subscribers map[string][]chan string
}

func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string][]chan string),
	}
}

func (ps *PubSub) Subscribe(topic string) chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan string)
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)

	return ch
}

func (ps *PubSub) Publish(topic, message string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, ch := range ps.subscribers[topic] {
		ch <- message
	}
}
