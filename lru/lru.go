package lru

import "container/list"

// cache for lru
// currently non-safe for concurrent
type Cache struct {
	// available used maximum bytes
	maxBytes int64
	// current used bytes
	nbytes int64
	// pointer for double linkedlist
	ll *list.List

	cache map[string]*list.Element

	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

// count the memory of return value
type Value interface {
	// method
	Len() int
}

// constructor of cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// get element by key and move it to front
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}
