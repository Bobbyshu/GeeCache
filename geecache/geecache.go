package geecache

import (
	"fmt"
	"geecache/singleflight"
	"log"
	"sync"
)

// A getter load data for key
// defines a method to load data for a given key.
type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function
type GetterFunc func(key string) ([]byte, error)

// Get calls f(key).
// This method implements the Get method of the Getter interface.
// When called, it invokes the function f with the given key.
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// A Group is a cache namespace and associated data loaded spread over
type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker
	// use singleflight.Group to make sure that
	// each key is only fetched once
	loader *singleflight.Group
}

var (
	// read-write mutex to protect the groups map.
	mu sync.RWMutex
	// map that stores all the created Group instances by their names.
	groups = make(map[string]*Group)
)

// Create a new instance of Group.
// Initialize a Group with the given name, cacheBytes, and getter.
// If the getter is nil, it panics.
// The newly created Group is stored in the global groups map.
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter!")
	}

	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}

	groups[name] = g
	return g
}

// Return the named group previously created with NewGroup, or
// nil if there's no such group.
// Use read lock (RLock) since it does not involve any write operations.
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	log.Printf("GetGroup called for '%s', groups: %v, result: %v", name, groups, g != nil)
	return g
}

// Get retrieves the value for a given key from the cache.
// It first checks the mainCache, and if the key is not found
// it calls the load method to fetch the data.
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	// If the key is not found in the cache, load the value.
	return g.load(key)
}

// load retrieves the value for a key, either locally or from a remote peer.
func (g *Group) load(key string) (value ByteView, err error) {
	if g.peers != nil {
		if peer, ok := g.peers.PickPeer(key); ok {
			if value, err = g.getFromPeer(peer, key); err == nil {
				return value, nil
			}
			log.Println("[GeeCache] Failed to get from peer", err)
		}
	}

	return g.getLocally(key)
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}

// getLocally loads the data for a key from the local source (getter).
func (g *Group) getLocally(key string) (ByteView, error) {
	// Call the user-defined getter to get the source data.
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	// Create a ByteView from the obtained bytes.
	value := ByteView{b: cloneBytes(bytes)}

	// Populate the cache with the obtained value.
	g.populateCache(key, value)

	return value, nil
}

// populateCache adds the key-value pair to the mainCache.
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

// RegisterPeers registers a PeerPicker for choosing remote peer
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}
