package consistent_hash

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type HashRing struct {
	mu       sync.RWMutex
	replicas int
	keys     []int
	hashMap  map[int]string
}

func New(replicas int) *HashRing {
	return &HashRing{
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
}

func (hr *HashRing) AddNode(node string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	for replicaIndex := 0; replicaIndex < hr.replicas; replicaIndex++ {
		hashKey := hr.hashKey(replicaIndex, node)
		hr.keys = append(hr.keys, hashKey)
		hr.hashMap[hashKey] = node
	}

	sort.Ints(hr.keys)
}

func (hr *HashRing) RemoveNode(node string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	for replicaIndex := 0; replicaIndex < hr.replicas; replicaIndex++ {
		hashKey := hr.hashKey(replicaIndex, node)
		index := sort.SearchInts(hr.keys, hashKey)

		keysBeforeIndex := hr.keys[:index]
		keysAfterIndex := hr.keys[index+1:]
		hr.keys = append(keysBeforeIndex, keysAfterIndex...)

		delete(hr.hashMap, hashKey)
	}
}

func (hr *HashRing) GetNode(node string) string {
	hr.mu.RLock()
	defer hr.mu.RUnlock()

	if len(hr.keys) == 0 {
		return ""
	}

	hash := hr.hashKey(0, node)
	index := sort.Search(len(hr.keys), func(i int) bool {
		return hr.keys[i] >= hash
	})

	if index == len(hr.keys) {
		index = 0
	}

	return hr.hashMap[hr.keys[index]]
}

func (hr *HashRing) hashKey(replicaIndex int, node string) int {
	key := strconv.Itoa(replicaIndex) + node
	return int(crc32.ChecksumIEEE([]byte(key)))
}
