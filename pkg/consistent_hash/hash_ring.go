package consistent_hash

import (
	"github.com/pedrokunz/distributed_cache_go/internal/data_node"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// HashRing is a consistent hash ring data structure
// - It manages the consistent hashing logic
// - It determines the node for a given key based on the hash
type HashRing struct {
	mu       sync.RWMutex
	replicas int
	keys     []int
	hashMap  map[int]*data_node.DataNode
}

func New(replicas int) *HashRing {
	return &HashRing{
		replicas: replicas,
		hashMap:  make(map[int]*data_node.DataNode),
	}
}

func (hr *HashRing) AddNode(node *data_node.DataNode) {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	for replicaIndex := 0; replicaIndex < hr.replicas; replicaIndex++ {
		hashKey := hr.hashKey(replicaIndex, node.ID())
		hr.keys = append(hr.keys, hashKey)
		hr.hashMap[hashKey] = node
	}

	sort.Ints(hr.keys)
}

func (hr *HashRing) RemoveNode(node *data_node.DataNode) {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	for replicaIndex := 0; replicaIndex < hr.replicas; replicaIndex++ {
		hashKey := hr.hashKey(replicaIndex, node.ID())
		index := sort.SearchInts(hr.keys, hashKey)

		keysBeforeIndex := hr.keys[:index]
		keysAfterIndex := hr.keys[index+1:]
		hr.keys = append(keysBeforeIndex, keysAfterIndex...)

		delete(hr.hashMap, hashKey)
	}
}

func (hr *HashRing) GetNode(key string) *data_node.DataNode {
	hr.mu.RLock()
	defer hr.mu.RUnlock()

	if len(hr.keys) == 0 {
		return nil
	}

	hash := hr.hashKey(0, key)
	index := sort.Search(len(hr.keys), func(i int) bool {
		return hr.keys[i] >= hash
	})

	if index == len(hr.keys) {
		index = 0
	}

	return hr.hashMap[hr.keys[index]]
}

func (hr *HashRing) Nodes() []*data_node.DataNode {
	hr.mu.RLock()
	defer hr.mu.RUnlock()

	nodes := make([]*data_node.DataNode, 0, len(hr.hashMap))
	for _, node := range hr.hashMap {
		nodes = append(nodes, node)
	}

	return nodes
}

func (hr *HashRing) hashKey(replicaIndex int, nodeID string) int {
	key := strconv.Itoa(replicaIndex) + nodeID
	return int(crc32.ChecksumIEEE([]byte(key)))
}
