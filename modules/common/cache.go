package common

import (
	"sync"
	"time"
)

type cacheEntry struct {
	results   []CategoryResult
	expiresAt time.Time
}

type Cache struct {
	entries sync.Map
	ttl     time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	if ttl <= 0 {
		ttl = 60 * time.Second
	}
	return &Cache{ttl: ttl}
}

func (c *Cache) Get(key string) ([]CategoryResult, bool) {
	v, ok := c.entries.Load(key)
	if !ok {
		return nil, false
	}
	entry := v.(cacheEntry)
	if time.Now().After(entry.expiresAt) {
		c.entries.Delete(key)
		return nil, false
	}
	return entry.results, true
}

func (c *Cache) Set(key string, results []CategoryResult) {
	c.entries.Store(key, cacheEntry{
		results:   results,
		expiresAt: time.Now().Add(c.ttl),
	})
}

type CachedFetcher struct {
	Source  string
	Fetcher Fetcher
	Cache   *Cache
}

func (c *CachedFetcher) FetchCategory(category string) ([]CategoryResult, error) {
	key := c.Source + ":" + category
	if results, ok := c.Cache.Get(key); ok {
		return results, nil
	}
	results, err := c.Fetcher.FetchCategory(category)
	if err != nil {
		return nil, err
	}
	if len(results) > 0 {
		c.Cache.Set(key, results)
	}
	return results, nil
}
