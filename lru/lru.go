package lru

import (
    "container/list"
)

/**
 * @author Rancho
 * @date 2021/12/15
 */

type Cache struct {
    // max allowed memory
    maxBytes int64
    // current used memory
    nbytes int64
    ll     *list.List
    cache  map[string]*list.Element
    // optional and executed when an entry is purged
    OnEvicted func(key string, value Value)
}

type entry struct {
    key   string
    value Value
}

type Value interface {
    Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
    return &Cache{
        maxBytes:  maxBytes,
        ll:        list.New(),
        cache:     make(map[string]*list.Element),
        OnEvicted: onEvicted,
    }
}

// Get look ups a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
    if item, ok := c.cache[key]; ok {
        c.ll.MoveToFront(item)
        kv := item.Value.(*entry)
        return kv.value, true
    }

    return
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
    item := c.ll.Back()
    if item == nil {
        return
    }

    c.ll.Remove(item)
    kv := item.Value.(*entry)
    delete(c.cache, kv.key)
    c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())

    if c.OnEvicted != nil {
        c.OnEvicted(kv.key, kv.value)
    }
}

func (c *Cache) Add(key string, value Value) {
    if item, ok := c.cache[key]; ok {
        c.ll.MoveToFront(item)
        kv := item.Value.(*entry)
        c.nbytes += int64(value.Len()) - int64(kv.value.Len())
        kv.value = value
    } else {
        item := c.ll.PushFront(&entry{
            key:   key,
            value: value,
        })
        c.cache[key] = item
        c.nbytes += int64(value.Len()) + int64(len(key))
    }

    for c.maxBytes != 0 && c.maxBytes < c.nbytes {
        c.RemoveOldest()
    }
}

// Size the number of cache entries
func (c *Cache) Size() int {
    return c.ll.Len()
}
