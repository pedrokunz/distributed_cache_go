package data_node

import (
	"log"
	"sync"
)

type DataNode struct {
	mu      sync.RWMutex
	storage map[string]string
}

func New() *DataNode {
	return &DataNode{
		storage: make(map[string]string),
	}
}

func (dn *DataNode) Get(key string) (string, bool) {
	dn.mu.RLock()
	defer dn.mu.RUnlock()

	value, exists := dn.storage[key]
	if !exists {
		log.Printf("Key %s does not exist\n", key)
	}

	return value, exists
}

func (dn *DataNode) Set(key, value string) {
	dn.mu.Lock()
	defer dn.mu.Unlock()

	dn.storage[key] = value
	log.Printf("Key %s is now set to %s\n", key, value)
}
