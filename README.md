# A distributed cache system

## Structure of code:
geecache/  
├── lru/  
│ └── lru.go // Least Recently Used caching strategy  
├── byteview.go // Encapsulation and Abstraction for cache  
├── cache.go // Concurrency Control  
└── geecache.go // External interaction (The main process of controlling cache storage and retrieval)  
