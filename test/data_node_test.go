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
