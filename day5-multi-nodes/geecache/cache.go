package geecache

import (
	"geecache/lru"
	"sync"
)

type cache struct {
	sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.Lock()
	defer c.Unlock()
	if c.lru == nil { // Lazy Initialization
		c.lru = lru.NewCache(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.Lock()
	defer c.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
