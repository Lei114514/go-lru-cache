package lru

import (
	"sync"
)

type ConcurrentCache struct{
	mu sync.RWMutex
	cache *LRUCache
}

func NewConcurrentCache(capacity int) *ConcurrentCache{
	return &ConcurrentCache{
		cache: NewLRUCache(capacity),
	}
}

func (c *ConcurrentCache) Get (key string) (value any, ok bool){
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache.Get(key)
}

func (c *ConcurrentCache) Put (key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Put(key,value)
}

func (c *ConcurrentCache) Delete(key string){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Delete(key)
}

func (c *ConcurrentCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Clear()
}

func (c *ConcurrentCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache.Len()
}

func (c *ConcurrentCache) Cap() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache.Cap()
}