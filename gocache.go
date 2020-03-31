package go_cache

import (
	"errors"
	"log"
	"sync"
)

var (
	rw     sync.RWMutex
	groups = make(map[string]*Group)
)

// Getter loads data for a key
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc is a callback function aim to fetch origin data and cache when data hadn't cached
type GetterFunc func(key string) ([]byte, error)

// transfer a external method to interface Getter
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group is a cache namespace and associated data loaded spread over
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

// NewGroup create a new instance of Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("getter is nil")
	}

	rw.Lock()
	defer rw.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}

	groups[name] = g
	return g
}

// GetGroup returns the named group previously created with NewGroup,
// or nil if there's no such group.
func GetGroup(name string) *Group {
	rw.RLock()
	defer rw.RUnlock()

	return groups[name]
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, errors.New("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GoCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

// load data from customer callback, Getter
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{data: byteClone(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// add a new data to cache
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
