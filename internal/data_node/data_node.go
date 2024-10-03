package data_node

import (
	"github.com/google/uuid"
	"github.com/pedrokunz/distributed_cache_go/internal/pub_sub"
	"log"
	"sync"
)

// DataNode represents a data node in the distributed cache system
type DataNode struct {
	mu    sync.RWMutex
	cache *LRUCache
	id    uuid.UUID
}

// New creates a new DataNode instance
func New() *DataNode {
	return &DataNode{
		cache: NewLRUCache(1000),
		id:    uuid.New(),
	}
}

// ID returns the ID of the data node
func (dn *DataNode) ID() string {
	return dn.id.String()
}

// Get retrieves the value for a given key from the cache
func (dn *DataNode) Get(key string) (string, bool) {
	dn.mu.RLock()
	defer dn.mu.RUnlock()

	value, exists := dn.cache.Get(key)
	if !exists {
		log.Printf("Key %s does not exist\n", key)
	}

	return value, exists
}

// Set sets the value for a given key in the cache
func (dn *DataNode) Set(key, value string) {
	dn.mu.Lock()
	defer dn.mu.Unlock()

	dn.cache.Set(key, value)

	log.Printf("Key %s is now set to %s\n", key, value)
}

// InvalidateCache invalidates a key in the cache
func (dn *DataNode) InvalidateCache(key string) {
	dn.mu.Lock()
	defer dn.mu.Unlock()

	dn.cache.Delete(key)

	log.Printf("Key %s has been invalidated\n", key)
}

// SubscribeToCacheInvalidation subscribes to cache invalidation events
func (dn *DataNode) SubscribeToCacheInvalidation(pubSub *pub_sub.PubSub) {
	invalidationChan := pubSub.Subscribe("cache_invalidation")
	go func() {
		for key := range invalidationChan {
			dn.InvalidateCache(key)
		}
	}()
}
