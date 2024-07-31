package cache

import (
    "time"
    "github.com/patrickmn/go-cache"
)

var c = cache.New(5*time.Minute, 10*time.Minute)

func Set(key string, value interface{}) {
    c.Set(key, value, cache.DefaultExpiration)
}

func Get(key string) (interface{}, bool) {
    return c.Get(key)
}
