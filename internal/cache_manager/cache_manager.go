package cache_manager

import (
	"log"
	"sync"
)

type CacheManager struct {
	mu        sync.RWMutex
	keyToNode map[string]string
}

func New() *CacheManager {
	return &CacheManager{
		keyToNode: make(map[string]string),
	}
}

func (cm *CacheManager) GetNodeForKey(key string) (string, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	node, exists := cm.keyToNode[key]

	return node, exists
}

func (cm *CacheManager) SetNodeForKey(key, node string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.keyToNode[key] = node

	log.Printf("Key %s is now on node %s\n", key, node)
}
