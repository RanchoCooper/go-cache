package lru

import (
	"container/list"
)

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

// Cache is a LRU cache. It's not safe for concurrent access.
type Cache struct {
	maxBytes   int64
	curBytes   int64
	linkedList *list.List
	cache      map[string]*list.Element

	// optional add executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:   maxBytes,
		linkedList: list.New(),
		cache:      make(map[string]*list.Element),
		OnEvicted:  onEvicted,
	}
}

// Get look ups a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.linkedList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	ele := c.linkedList.Back()
	if ele != nil {
		c.linkedList.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.curBytes -= int64(len(kv.key)) + int64(kv.value.Len())

		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add adds a value to the cache.
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// update with new value
		c.linkedList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.curBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.linkedList.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.cache[key] = ele
		c.curBytes += int64(len(key)) + int64(value.Len())
	}

	// scale cache size if overload
	for c.maxBytes != 0 && c.maxBytes < c.curBytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.linkedList.Len()
}

