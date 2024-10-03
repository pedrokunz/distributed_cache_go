package test

import (
	"strconv"
	"sync"
	"testing"

	"github.com/pedrokunz/distributed_cache_go/internal/cache_manager"
)

func TestCacheManager(t *testing.T) {
	cm := cache_manager.New()

	// Test setting a new key-value pairs
	err := cm.SetNodeForKey("key1", "value1")
	if err != nil {
		t.Errorf("error setting key-value pair: %v", err)
	}

	node, err := cm.GetNodeForKey("key1")
	if err != nil {
		t.Errorf("error getting key-value pair: %v", err)
	}

	if node == nil {
		t.Fatal("expected node for key1, got none")
	}

	value, _ := node.Get("key1")
	if value != "value1" {
		t.Errorf("expected value1 for key1, got %s", value)
	}

	// Test updating key-value pairs
	err = cm.SetNodeForKey("key1", "value2")
	if err != nil {
		t.Errorf("error setting key-value pair: %v", err)
	}

	node, err = cm.GetNodeForKey("key1")
	if err != nil {
		t.Errorf("error getting key-value pair: %v", err)
	}

	if node == nil {
		t.Fatal("expected node for key1, got none")
	}

	value, _ = node.Get("key1")
	if value != "value2" {
		t.Errorf("expected value2 for key1, got %s", value)
	}

	// Test deleting key-value pairs
	err = cm.DeleteNodeForKey("key1")
	if err != nil {
		t.Errorf("error deleting key-value pair: %v", err)
	}
	node, err = cm.GetNodeForKey("key1")
	if err == nil {
		t.Errorf("expected no node for key1, got one %v", node)
	}
}

func TestCacheManagerConcurrent(t *testing.T) {
	cm := cache_manager.New()

	var wg sync.WaitGroup
	const numGoroutines = 100

	// Test concurrent setting of key-value pairs
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			value := "value" + strconv.Itoa(i)
			err := cm.SetNodeForKey(key, value)
			if err != nil {
				t.Errorf("error setting key-value pair: %v", err)
			}
		}(i)
	}
	wg.Wait()

	// Test concurrent getting of key-value pairs
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			_, err := cm.GetNodeForKey(key)
			if err != nil {
				t.Errorf("error getting key-value pair: %v", err)
			}
		}(i)
	}
	wg.Wait()

	// Test concurrent deleting of key-value pairs
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			err := cm.DeleteNodeForKey(key)
			if err != nil {
				t.Errorf("error deleting key-value pair: %v", err)
			}
		}(i)
	}

	wg.Wait()
}
