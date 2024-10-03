package data_node

import (
	"github.com/google/uuid"
	"github.com/pedrokunz/distributed_cache_go/internal/pub_sub"
	"log"
	"sync"
)

type DataNode struct {
	mu    sync.RWMutex
	cache *LRUCache
	id    uuid.UUID
}

func New() *DataNode {
	return &DataNode{
		cache: NewLRUCache(1000),
		id:    uuid.New(),
	}
}

func (dn *DataNode) ID() string {
	return dn.id.String()
}

func (dn *DataNode) Get(key string) (string, bool) {
	dn.mu.RLock()
	defer dn.mu.RUnlock()

	value, exists := dn.cache.Get(key)
	if !exists {
		log.Printf("Key %s does not exist\n", key)
	}

	return value, exists
}

func (dn *DataNode) Set(key, value string) {
	dn.mu.Lock()
	defer dn.mu.Unlock()

	dn.cache.Set(key, value)

	log.Printf("Key %s is now set to %s\n", key, value)
}

func (dn *DataNode) InvalidateCache(key string) {
	dn.mu.Lock()
	defer dn.mu.Unlock()

	dn.cache.Delete(key)

	log.Printf("Key %s has been invalidated\n", key)
}

func (dn *DataNode) SubscribeToCacheInvalidation(pubSub *pub_sub.PubSub) {
	invalidationChan := pubSub.Subscribe("cache_invalidation")
	go func() {
		for key := range invalidationChan {
			dn.InvalidateCache(key)
		}
	}()
}
