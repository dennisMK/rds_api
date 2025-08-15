package concurrent

import (
	"sync"
	"time"
)

// CacheItem represents a cached item with expiration
type CacheItem[T any] struct {
	Value     T
	ExpiresAt time.Time
}

// IsExpired checks if the cache item has expired
func (ci *CacheItem[T]) IsExpired() bool {
	return time.Now().After(ci.ExpiresAt)
}

// ConcurrentCache provides thread-safe caching with TTL
type ConcurrentCache[K comparable, V any] struct {
	items map[K]*CacheItem[V]
	mutex sync.RWMutex
	ttl   time.Duration
}

// NewConcurrentCache creates a new concurrent cache
func NewConcurrentCache[K comparable, V any](ttl time.Duration) *ConcurrentCache[K, V] {
	cache := &ConcurrentCache[K, V]{
		items: make(map[K]*CacheItem[V]),
		ttl:   ttl,
	}
	
	// Start cleanup goroutine
	go cache.cleanup()
	
	return cache
}

// Set stores a value in the cache
func (c *ConcurrentCache[K, V]) Set(key K, value V) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.items[key] = &CacheItem[V]{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Get retrieves a value from the cache
func (c *ConcurrentCache[K, V]) Get(key K) (V, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	item, exists := c.items[key]
	if !exists || item.IsExpired() {
		var zero V
		return zero, false
	}
	
	return item.Value, true
}

// Delete removes a value from the cache
func (c *ConcurrentCache[K, V]) Delete(key K) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	delete(c.items, key)
}

// Clear removes all items from the cache
func (c *ConcurrentCache[K, V]) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.items = make(map[K]*CacheItem[V])
}

// Size returns the number of items in the cache
func (c *ConcurrentCache[K, V]) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	return len(c.items)
}

// cleanup removes expired items from the cache
func (c *ConcurrentCache[K, V]) cleanup() {
	ticker := time.NewTicker(c.ttl / 2)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mutex.Lock()
		for key, item := range c.items {
			if item.IsExpired() {
				delete(c.items, key)
			}
		}
		c.mutex.Unlock()
	}
}
