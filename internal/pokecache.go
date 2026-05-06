// Package responsible for handling the caching of responses from the PokeAPI.
package pokecache

import (
	"log"
	"sync"
	"time"
)

type Cache struct {
	entries map[string]CacheEntry
	mutex   sync.RWMutex
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entries: map[string]CacheEntry{},
	}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, ok := c.entries[key]
	if ok {
		// entry already exists
		return
	}
	c.entries[key] = CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, ok := c.entries[key]
	if ok {
		return entry.val, true
	}
	return nil, false
}

// Delete all cache entries that are older than the specified time.
// This prevents the cache from growing too large over time.
// This function is meant to start running asynchronously after a new Cache object is created.
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for tickTime := range ticker.C {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		toDelete := []string{}
		for key, entry := range c.entries {
			if tickTime.Sub(entry.createdAt) >= interval {
				toDelete = append(toDelete, key)
				log.Printf("found expired cache entry: %s", key)
			}
		}

		for _, key := range toDelete {
			delete(c.entries, key)
		}
	}
}
