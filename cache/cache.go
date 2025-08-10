package cache

import (
	"fmt"
	"sync"
)

// Cache represents the hierarchical cache structure
type Cache struct {
	mu       sync.RWMutex
	data     map[string][]byte // key -> data mapping
	resolver *KeyResolver
}

// NewCache creates a new cache instance with the given configuration
func NewCache(config *TableConfig) *Cache {
	return &Cache{
		data:     make(map[string][]byte),
		resolver: NewKeyResolver(config),
	}
}

// Get retrieves data from cache by key
func (c *Cache) Get(key string) ([]byte, bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	// Validate key format
	if !c.resolver.ValidateKey(key) {
		return nil, false, fmt.Errorf("invalid key format: %s", key)
	}
	
	data, exists := c.data[key]
	if !exists {
		return nil, false, nil
	}
	
	// Return a copy to prevent external modification
	result := make([]byte, len(data))
	copy(result, data)
	
	return result, true, nil
}

// Set stores data in cache with the specified key
func (c *Cache) Set(key string, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Validate key format
	if !c.resolver.ValidateKey(key) {
		return fmt.Errorf("invalid key format: %s", key)
	}
	
	// Store a copy to prevent external modification
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	
	c.data[key] = dataCopy
	
	return nil
}

// Delete removes data from cache by key
func (c *Cache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Validate key format
	if !c.resolver.ValidateKey(key) {
		return fmt.Errorf("invalid key format: %s", key)
	}
	
	delete(c.data, key)
	
	return nil
}

// Clear removes all data from cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.data = make(map[string][]byte)
}

// Size returns the number of items in cache
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return len(c.data)
}

// Keys returns all keys in cache
func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	keys := make([]string, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}
	
	return keys
}