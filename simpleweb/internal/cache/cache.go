package cache

import (
    "sync"
)

var (
    cacheMap map[string]interface{}
    mutex    sync.Mutex
)

func init() {
    cacheMap = make(map[string]interface{})
}

func Set(key string, value interface{}) {
    mutex.Lock()
    defer mutex.Unlock()
    cacheMap[key] = value
}

func Get(key string) (interface{}, bool) {
    mutex.Lock()
    defer mutex.Unlock()
    value, found := cacheMap[key]
    return value, found
}
