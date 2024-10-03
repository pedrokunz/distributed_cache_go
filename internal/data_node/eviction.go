package data_node

import (
	"container/list"
	"log"
)

type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
}

type entry struct {
	key   string
	value string
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (c *LRUCache) Get(key string) (string, bool) {
	elem, exists := c.cache[key]
	if exists {
		c.list.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}

	return "", false
}

func (c *LRUCache) Set(key, value string) {
	elem, exists := c.cache[key]
	if exists {
		c.list.MoveToFront(elem)
		elem.Value.(*entry).value = value
	} else {
		if c.list.Len() == c.capacity {
			c.evict()
		}

		el := c.list.PushFront(&entry{key, value})
		c.cache[key] = el
	}

	log.Printf("Key %s is now set to %s\n", key, value)
}

func (c *LRUCache) evict() {
	el := c.list.Back()
	if el != nil {
		c.list.Remove(el)

		kv := el.Value.(*entry)
		delete(c.cache, kv.key)

		log.Printf("Key %s has been evicted\n", kv.key)
	}
}
