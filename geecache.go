package geecache

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
