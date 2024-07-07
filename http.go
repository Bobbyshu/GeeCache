package geecache

import (
	"fmt"
	"log"
)

const defaultBasePath = "/_geecache"

// Implements PeerPicker for a pool of HTTP peers.
type HTTPPool struct {
	// register self address localhost(ip):port
	self     string
	basePath string
}

// Initializes an HTTP pool of peers.
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log info with server name
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}
