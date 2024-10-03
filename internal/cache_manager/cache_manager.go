package cache_manager

import (
	"fmt"
	"github.com/pedrokunz/distributed_cache_go/internal/data_node"
	"github.com/pedrokunz/distributed_cache_go/internal/pub_sub"
	"github.com/pedrokunz/distributed_cache_go/pkg/consistent_hash"
	"log"
	"sync"
)

// CacheManager is the main struct for the cache manager
// - It manages the mapping of keys to nodes
// - It uses the HashRing to determine the appropriate node for a given key
// - It uses the PubSub to publish node key assignments and cache invalidations
type CacheManager struct {
	mu        sync.RWMutex
	keyToNode map[string]*data_node.DataNode
	pubSub    *pub_sub.PubSub
	hashRing  *consistent_hash.HashRing
	nodes     map[string]*data_node.DataNode
}

func New() *CacheManager {
	return &CacheManager{
		keyToNode: make(map[string]*data_node.DataNode),
		pubSub:    pub_sub.NewPubSub(),
		hashRing:  consistent_hash.New(3),
		nodes:     make(map[string]*data_node.DataNode),
	}
}

func (cm *CacheManager) GetNodeForKey(key string) (*data_node.DataNode, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	node := cm.hashRing.GetNode(key)
	if node == nil {
		log.Printf("No node found for key %s\n", key)
		return nil, fmt.Errorf("no node found for key %s", key)
	}

	return node, nil
}

func (cm *CacheManager) SetNodeForKey(key, value string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	node := cm.hashRing.GetNode(key)
	if node == nil {
		node = data_node.New()
		cm.hashRing.AddNode(node)
		cm.nodes[node.ID()] = node
	}

	cm.keyToNode[key] = node
	node.Set(key, value)

	log.Printf("Key %s is now on node %s with value %s\n", key, node.ID(), value)

	cm.pubSub.Publish("node_key_assignment", key)
	cm.pubSub.Publish("cache_invalidation", key)

	return nil
}

func (cm *CacheManager) DeleteNodeForKey(key string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	node := cm.keyToNode[key]
	if node == nil {
		log.Printf("No node found for key %s\n", key)
		return fmt.Errorf("no node found for key %s", key)
	}

	node.InvalidateCache(key)
	delete(cm.keyToNode, key)

	hasKeys := false
	for _, n := range cm.keyToNode {
		if n == node {
			hasKeys = true
			break
		}
	}

	// Remove the node from the hash ring if it has no remaining keys
	if !hasKeys {
		cm.hashRing.RemoveNode(node)
		delete(cm.nodes, node.ID())
	}

	log.Printf("Key %s has been deleted from node %s\n", key, node.ID())

	cm.pubSub.Publish("cache_invalidation", key)

	return nil
}

func (cm *CacheManager) Nodes() []*data_node.DataNode {
	return cm.hashRing.Nodes()
}

func (cm *CacheManager) SubscribeToNodeKeyAssignment() chan string {
	return cm.pubSub.Subscribe("node_key_assignment")
}

func (cm *CacheManager) SubscribeToCacheInvalidation() chan string {
	return cm.pubSub.Subscribe("cache_invalidation")
}
