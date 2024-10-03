package consistent_hash

import "sync"

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
