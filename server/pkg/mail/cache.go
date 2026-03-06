package mail

import (
	"sync"
	"time"
)

const cacheTTL = 10 * time.Minute

// cacheEntry holds a cached rendered template
type cacheEntry struct {
	body      string
	expiresAt time.Time
}

// renderCache is a thread-safe TTL cache for rendered email templates.
// The cache key is the template name + a hash of the data.
type renderCache struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry
}

func newRenderCache() *renderCache {
	c := &renderCache{
		entries: make(map[string]cacheEntry),
	}
	go c.evictLoop()
	return c
}

func (c *renderCache) get(key string) (string, bool) {
	c.mu.RLock()
	entry, ok := c.entries[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(entry.expiresAt) {
		return "", false
	}
	return entry.body, true
}

func (c *renderCache) set(key, body string) {
	c.mu.Lock()
	c.entries[key] = cacheEntry{
		body:      body,
		expiresAt: time.Now().Add(cacheTTL),
	}
	c.mu.Unlock()
}

func (c *renderCache) evictLoop() {
	ticker := time.NewTicker(cacheTTL)
	defer ticker.Stop()
	for range ticker.C {
		c.evict()
	}
}

func (c *renderCache) evict() {
	now := time.Now()
	c.mu.Lock()
	for key, entry := range c.entries {
		if now.After(entry.expiresAt) {
			delete(c.entries, key)
		}
	}
	c.mu.Unlock()
}
