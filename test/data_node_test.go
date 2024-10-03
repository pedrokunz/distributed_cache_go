package test

import (
	"strconv"
	"sync"
	"testing"

	"github.com/pedrokunz/distributed_cache_go/internal/data_node"
)

func TestDataNode_SetAndGet(t *testing.T) {
	node := data_node.New()

	// Test setting and getting a key-value pair
	node.Set("key1", "value1")
	value, exists := node.Get("key1")
	if !exists || value != "value1" {
		t.Errorf("expected value1 for key1, got %s", value)
	}
}

func TestDataNode_InvalidateCache(t *testing.T) {
	node := data_node.New()

	// Test invalidating a key
	node.Set("key1", "value1")
	node.InvalidateCache("key1")
	_, exists := node.Get("key1")
	if exists {
		t.Errorf("expected key1 to be invalidated")
	}
}

func TestDataNode_ConcurrentSetAndGet(t *testing.T) {
	node := data_node.New()
	var wg sync.WaitGroup
	const numGoroutines = 100

	// Test concurrent setting of key-value pairs
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			value := "value" + strconv.Itoa(i)
			node.Set(key, value)
		}(i)
	}
	wg.Wait()

	// Test concurrent getting of key-value pairs
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			value, exists := node.Get(key)
			if !exists || value != "value"+strconv.Itoa(i) {
				t.Errorf("expected value%s for key%s, got %s", strconv.Itoa(i), strconv.Itoa(i), value)
			}
		}(i)
	}
	wg.Wait()
}

func TestLRUCacheEviction(t *testing.T) {
	capacity := 100
	lru := data_node.NewLRUCache(capacity)

	// Insert more items than the capacity to trigger eviction
	for i := 0; i < capacity*2; i++ {
		key := "key" + strconv.Itoa(i)
		value := "value" + strconv.Itoa(i)
		lru.Set(key, value)
	}

	// Check that the cache size does not exceed the capacity
	if len(lru.Cache) > capacity {
		t.Errorf("expected cache size to be %d, but got %d", capacity, len(lru.Cache))
	}

	// Check that the first inserted items have been evicted
	for i := 0; i < capacity; i++ {
		key := "key" + strconv.Itoa(i)
		if _, exists := lru.Get(key); exists {
			t.Errorf("expected key %s to be evicted, but it still exists", key)
		}
	}

	// Check that the last inserted items are still in the cache
	for i := capacity; i < capacity*2; i++ {
		key := "key" + strconv.Itoa(i)
		if value, exists := lru.Get(key); !exists || value != "value"+strconv.Itoa(i) {
			t.Errorf("expected key %s to be in the cache with value %s, but got %s", key, "value"+strconv.Itoa(i), value)
		}
	}
}
