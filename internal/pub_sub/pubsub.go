package pub_sub

import "sync"

// PubSub is a simple pub-sub implementation
type PubSub struct {
	mu          sync.RWMutex
	subscribers map[string][]chan string
}

// New creates a new PubSub instance
func New() *PubSub {
	return &PubSub{
		subscribers: make(map[string][]chan string),
	}
}

// Subscribe subscribes to a topic
func (ps *PubSub) Subscribe(topic string) chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan string)
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)

	return ch
}

// Publish publishes a message to a topic
func (ps *PubSub) Publish(topic, message string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, ch := range ps.subscribers[topic] {
		ch <- message
	}
}
