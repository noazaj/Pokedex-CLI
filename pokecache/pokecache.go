package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	data            map[string]cacheEntry
	defaultDuration time.Duration
	mu              *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

var (
	globalCache *Cache
	once        sync.Once
)

// Create a globalCache instance for functions to use
func InitGlobalCache() {
	once.Do(func() {
		// Initialize with a default duration or a configurable one
		globalCache = NewCache(5 * time.Minute)
	})
}

// GetGlobalCache returns the singleton instance of the cache
func GetGlobalCache() *Cache {
	return globalCache
}

// This function creates a new cache for the request data
func NewCache(interval time.Duration) *Cache {
	newCache := Cache{
		data:            make(map[string]cacheEntry),
		defaultDuration: interval,
		mu:              new(sync.Mutex),
	}
	go newCache.reapLoop()
	return &newCache
}

// Adds a new cached entry to the cache with the key as the url and
// the data as the cacheEntry
func (c *Cache) Add(key string, val []byte) {
	// Lock and then defer the unlock of the map for concurrency issues
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create a new cache entry
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}

	// Add the entry
	c.data[key] = entry
}

// Gets a specific field of data from the given string if present in the cache
func (c *Cache) Get(key string) ([]byte, bool) {
	// Lock and unlock for concurrency purposes
	c.mu.Lock()
	defer c.mu.Unlock()

	// Read from the cache and see if the string is present
	value, ok := c.data[key]
	if !ok {
		return nil, false
	}
	return value.val, true
}

// Goroutine to run concurrently to check the timing of the cached values
// If the values are old and greater than the time duration, delete them
// from the cache
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.defaultDuration)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				c.mu.Lock()
				for k, v := range c.data {
					if time.Since(v.createdAt) > c.defaultDuration {
						delete(c.data, k)
					}
				}
				c.mu.Unlock()
			}
		}
	}()
}
