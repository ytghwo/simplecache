package simplecache

import (
	"sync"

	"github.com/ytghwo/simplecache/lru"
)

//控制并发读写，可切换缓存淘汰逻辑

type cache struct {
	mu       sync.Mutex
	lru      *lru.Cache
	capacity int
}

func NewCache(capacity int) *cache {
	return &cache{capacity: capacity}
}

func (c *cache) add(key string, value byteview) {
	if c.lru == nil {
		c.lru = lru.New(int64(c.capacity), nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (byteview, bool) {
	if c.lru == nil {
		return byteview{}, false
	}
	c.mu.Lock()
	if v, ok := c.lru.Get(key); ok {
		return v.(byteview), true
	}
	return byteview{}, false
}
