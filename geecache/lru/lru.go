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
		// update ele
		c.ll.MoveToFront(ele)
		// predicate in Go x.(T)
		// cache map[string]*list.Element
		// string -> pointer
		kv := ele.Value.(*entry)
		return kv.value, true
	}

	// can't find
	return
}

// remove least recently used element
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		// delete the key store in map
		delete(c.cache, kv.key)

		// update used bytes in cache
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// update corresponding kv
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	// if key exist
	if ele, ok := c.cache[key]; ok {
		// move it to front and update value & bytes
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// non-exist -> add new entry
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	// remove element over maxBytes
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// check the amount of entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
