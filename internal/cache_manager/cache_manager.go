package cache_manager

import (
	"log"
	"sync"
)

type CacheManager struct {
	mu        sync.RWMutex
	keyToNode map[string]string
	pubSub    *PubSub
}

func New() *CacheManager {
	return &CacheManager{
		keyToNode: make(map[string]string),
		pubSub:    NewPubSub(),
	}
}

func (cm *CacheManager) GetNodeForKey(key string) (string, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	node, exists := cm.keyToNode[key]
	if !exists {
		log.Printf("Key %s does not exist\n", key)
	}

	return node, exists
}

func (cm *CacheManager) SetNodeForKey(key, node string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.keyToNode[key] = node

	log.Printf("Key %s is now on node %s\n", key, node)

	cm.pubSub.Publish("node_key_assignment", key)
}

func (cm *CacheManager) SubscribeToNodeKeyAssignment() chan string {
	return cm.pubSub.Subscribe("node_key_assignment")
}
