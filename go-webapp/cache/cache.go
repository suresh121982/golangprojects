// cache.go
package cache

import (
    "sync"
)

var (
    cacheMap = make(map[string]interface{})
    mu       sync.RWMutex
)

// SetCache sets a key-value pair in cache
func SetCache(key string, value interface{}) {
    mu.Lock()
    defer mu.Unlock()
    cacheMap[key] = value
}

// GetCache retrieves a value from cache by key
func GetCache(key string) (interface{}, bool) {
    mu.RLock()
    defer mu.RUnlock()
    val, ok := cacheMap[key]
    return val, ok
}

// DeleteCache deletes a value from cache by key
func DeleteCache(key string) {
    mu.Lock()
    defer mu.Unlock()
    delete(cacheMap, key)
}
