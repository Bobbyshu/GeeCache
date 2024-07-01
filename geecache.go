package geecache

import "sync"

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
		name: name,
		getter: getter,
		mainCache: cache{cacheBytes: cacheBytes}
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
	return g
}
