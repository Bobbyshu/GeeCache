package geecache

import (
	"geecache/lru"
	"sync"
)

type cache struct {
	// make sure only one coroutine can get and update cache
	mu sync.Mutex

	lru *lru.Cache

	// maximum cap
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	// release lock after method
	defer c.mu.Lock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	// add cache
	c.lru.Add(key, value)
}
