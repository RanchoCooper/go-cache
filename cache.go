package go_cache

import (
	"sync"

	"github.com/RanchoCooper/go-cache/lru"
)

type cache struct {
	rw sync.RWMutex
	lru *lru.Cache
	cacheBytes int64
}

// add will initialize lru if necessary
func (c *cache) add(key string, value ByteView) {
	c.rw.Lock()
	defer c.rw.Unlock()

	if c.lru == nil {
		// lazy initialization
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
