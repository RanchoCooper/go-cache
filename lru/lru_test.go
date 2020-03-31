package lru

import (
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestCache_Get(t *testing.T) {
	lru := New(int64(0), nil)

	key := "key1"
	value := String("1234")
	lru.Add(key, value)
	if v, ok := lru.Get(key); !ok {
		t.Fatalf("cache hit key faild.")
	} else if  v != value {
		t.Fatalf("cache value not equal.")
	}

	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache got success with non-exists key")
	}
}