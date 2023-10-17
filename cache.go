package main

import (
	"sync"
	"time"
)

type CacheItem struct {
	key     string
	value   []byte
	expires time.Time
	next    *CacheItem
	prev    *CacheItem
}

type Cache struct {
	mu   sync.RWMutex
	head *CacheItem
	tail *CacheItem
	stop chan struct{}
}

const CacheCapacity = 1000

func NewCache() *Cache {
	cache := &Cache{
		stop: make(chan struct{}),
	}
	go cache.startCleanup()
	return cache
}

func (c *Cache) Set(key []byte, value []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	keyStr := string(key)
	currentTime := time.Now()
	expirationTime := currentTime.Add(ttl)

	item := &CacheItem{
		key:     keyStr,
		value:   value,
		expires: expirationTime,
	}

	if c.head == nil {
		c.head = item
		c.tail = item
	} else {
		item.next = c.head
		c.head.prev = item
		c.head = item

		if count := c.countItems(); count > CacheCapacity {
			c.removeOldestExpired()
		}
	}
}

func (c *Cache) Get(key []byte) (value []byte, ttl time.Duration) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keyStr := string(key)
	currentTime := time.Now()
	for item := c.head; item != nil; item = item.next {
		if item.key == keyStr {
			if item.expires.After(currentTime) {
				return item.value, item.expires.Sub(currentTime)
			}
			// Item has expired, remove it
			c.removeCacheItem(item)
		}
	}

	return nil, 0
}

func (c *Cache) startCleanup() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c.stop:
			return
		case <-ticker.C:
			c.mu.Lock()
			c.removeExpired()
			c.mu.Unlock()
		}
	}
}

func (c *Cache) removeExpired() {
	currentTime := time.Now()
	for item := c.head; item != nil; {
		if item.expires.Before(currentTime) {
			nextItem := item.next
			c.removeCacheItem(item)
			item = nextItem
		} else {
			break
		}
	}
}

func (c *Cache) removeOldestExpired() {
	for item := c.tail; item != nil; item = item.prev {
		if item.expires.Before(time.Now()) {
			c.removeCacheItem(item)
			break
		}
	}
}

func (c *Cache) removeCacheItem(item *CacheItem) {
	if item == c.head {
		c.head = item.next
	}
	if item == c.tail {
		c.tail = item.prev
	}
	if item.next != nil {
		item.next.prev = item.prev
	}
	if item.prev != nil {
		item.prev.next = item.next
	}
}

func (c *Cache) StopCleanup() {
	close(c.stop)
}

func (c *Cache) countItems() int {
	count := 0
	for item := c.head; item != nil; item = item.next {
		count++
	}
	return count
}

func main() {
	cache := NewCache()
	defer cache.StopCleanup()

	cache.Set([]byte("DJ"), []byte("Dante J"), 5*time.Second)
	cache.Set([]byte("DM"), []byte("MOtley"), 10*time.Second)

	value, ttl := cache.Get([]byte("DJ"))
	println("Value:", string(value), "\t TTL:", ttl)

	value, ttl = cache.Get([]byte("DM"))
	println("Value:", string(value), "\t TTL:", ttl)
}
