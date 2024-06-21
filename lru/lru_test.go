package lru

import (
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)

	lru.Add("key1", String("1003"))

	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1003" {
		t.Fatalf("cache hit key1=1003 failed")
	}

	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}
