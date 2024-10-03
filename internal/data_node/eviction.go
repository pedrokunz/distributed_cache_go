package data_node

import (
	"container/list"
	"log"
)

type LRUCache struct {
	capacity   int
	Cache      map[string]*list.Element
	list       *list.List
	usageCount map[string]int
}

type entry struct {
	key   string
	value string
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity:   capacity,
		Cache:      make(map[string]*list.Element),
		list:       list.New(),
		usageCount: make(map[string]int),
	}
}

func (c *LRUCache) Get(key string) (string, bool) {
	elem, exists := c.Cache[key]
	if exists {
		c.list.MoveToFront(elem)
		c.usageCount[key]++
		return elem.Value.(*entry).value, true
	}

	return "", false
}

func (c *LRUCache) Set(key, value string) {
	elem, exists := c.Cache[key]
	if exists {
		c.list.MoveToFront(elem)
		elem.Value.(*entry).value = value
	} else {
		if c.list.Len() == c.capacity {
			c.evict()
		}

		el := c.list.PushFront(&entry{key, value})
		c.Cache[key] = el
		c.usageCount[key] = 0
	}

	c.usageCount[key]++

	log.Printf("Cache: Key %s is now set to %s\n", key, value)
}

func (c *LRUCache) Delete(key string) {
	elem, exists := c.Cache[key]
	if exists {
		c.list.Remove(elem)

		delete(c.Cache, key)
		delete(c.usageCount, key)

		log.Printf("Cache: Key %s has been deleted\n", key)
	}
}

func (c *LRUCache) evict() {
	el := c.list.Back()
	if el != nil {
		c.list.Remove(el)

		kv := el.Value.(*entry)
		delete(c.Cache, kv.key)
		delete(c.usageCount, kv.key)

		log.Printf("Cache: Key %s has been evicted\n", kv.key)
	}
}
