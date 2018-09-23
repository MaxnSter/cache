package simple

import (
	"fmt"
	"github.com/MaxnSter/cache"
	"github.com/jonboulle/clockwork"
	"golang.org/x/sync/singleflight"
	"sync"
	"time"
)

const cacheName = "simple"

type cacheItem struct {
	clock      clockwork.Clock
	value      cache.Value
	expiration *time.Time
}

type simpleCache struct {
	*cache.Stats
	l sync.Mutex // guard cache

	cache  map[cache.Key]*cacheItem
	config cache.Config
	clock  clockwork.Clock
	group  *singleflight.Group
}

func init() {
	cache.RegisterCreator(cacheName, newSimpleCache)
}

func newSimpleCache(config *cache.Config) cache.Cache {
	c := &simpleCache{
		l:      sync.Mutex{},
		config: *config,
		clock:  clockwork.NewRealClock(),
		group:  &singleflight.Group{},
		Stats:  &cache.Stats{},
	}

	c.init()
	return c
}

func (c *simpleCache) init() {
	if c.config.Size > 0 {
		c.cache = make(map[cache.Key]*cacheItem, c.config.Size)
	} else {
		c.cache = map[cache.Key]*cacheItem{}
	}
}

func (c *simpleCache) Set(k cache.Key, v cache.Value) error {
	c.l.Lock()
	defer c.l.Unlock()

	return c.setLocked(k, v)
}

func (c *simpleCache) setLocked(k cache.Key, v cache.Value) error {
	var err error
	if c.config.Serializer != nil {
		v, err = c.config.Serializer.Serialize(k, v)
		if err != nil {
			return err
		}
	}

	item, ok := c.cache[k]
	if ok {
		item.value = v
	} else {

		if c.config.Size > 0 && len(c.cache) >= c.config.Size {
			c.evictLocked(1)
		}

		item = &cacheItem{
			clock: c.clock,
			value: v,
		}
		c.cache[k] = item
	}

	if c.config.Expiration != nil {
		t := item.expiration.Add(*c.config.Expiration)
		item.expiration = &t
	}

	// FIXME callback block or deadlock
	if c.config.AddedFunc != nil {
		c.config.AddedFunc(k, v)
	}
	return nil
}

func (c *simpleCache) evictLocked(count int) {
	now := c.clock.Now()
	var evicted int

	for k, item := range c.cache {
		if evicted >= count {
			return
		}
		if item.expiration == nil || now.After(*item.expiration) {
			c.removeLocked(k)
			evicted++
		}
	}
}

func (c *simpleCache) SetWithExpire(k cache.Key, v cache.Value, t time.Duration) error {
	c.l.Lock()
	defer c.l.Unlock()

	if err := c.setLocked(k, v); err != nil {
		return err
	}

	nt := c.clock.Now().Add(t)
	c.cache[k].expiration = &nt
	return nil
}

func (c *simpleCache) Get(k cache.Key) (cache.Value, bool) {
	item, ok := c.lookup(k)
	if !ok {
		c.IncrMissCount()
		return nil, false
	}

	var v cache.Value
	var err error
	if c.config.Deserializer != nil {
		v, err = c.config.Deserializer.Deserialize(k, item.value)
		if err != nil {
			return nil, false
		}
	}

	c.IncrHitCount()
	return v, true
}

func (c *simpleCache) lookup(k cache.Key) (*cacheItem, bool) {
	c.l.Lock()
	defer c.l.Unlock()

	item, ok := c.cache[k]
	if !ok {
		return nil, false
	}

	if item.expiration != nil &&
		c.clock.Now().After(*item.expiration) {
		c.removeLocked(k)
		return nil, false
	}

	return item, true
}

func (c *simpleCache) GetALL() map[cache.Key]cache.Value {
	m := map[cache.Key]cache.Value{}
	for k := range c.cache {
		item, ok := c.lookup(k)

		if !ok {
			continue
		}
		// FIXME  value deserialize
		m[k] = item.value
	}
	return m
}

func (c *simpleCache) Load(k cache.Key) error {
	if c.config.Loader == nil {
		return fmt.Errorf("loader not found in cacheConfig")
	}

	return c.load(k, func(k cache.Key, v cache.Value) error {
		return c.Set(k, v)
	})
}

func (c *simpleCache) LoadWithExpire(k cache.Key, t time.Duration) error {
	if c.config.Loader == nil {
		return fmt.Errorf("loader not found in cacheConfig")
	}

	return c.load(k, func(k cache.Key, v cache.Value) error {
		return c.SetWithExpire(k, v, t)
	})
}

func (c *simpleCache) load(k cache.Key, afterLoaded func(k cache.Key, v cache.Value) error) error {
	_, err, _ := c.group.Do(k, func() (interface{}, error) {
		if item, ok := c.lookup(k); ok {
			return item.value, nil
		}

		v, err := c.config.Loader.Load(k)
		if err != nil {
			return nil, err
		}

		if afterLoaded != nil {
			if err := afterLoaded(k, v); err != nil {
				return nil, err
			}
		}
		return v, nil
	})

	return err
}

func (c *simpleCache) Remove(k cache.Key) bool {
	c.l.Lock()
	defer c.l.Unlock()

	return c.removeLocked(k)
}

func (c *simpleCache) removeLocked(k cache.Key) bool {
	item, ok := c.cache[k]
	if !ok {
		return false
	}

	delete(c.cache, k)
	if c.config.EvictedFunc != nil {
		c.config.EvictedFunc(k, item.value)
	}
	return true
}

func (c *simpleCache) Purge() {
	c.l.Lock()
	defer c.l.Unlock()

	if c.config.EvictedFunc != nil {
		for k, item := range c.cache {
			c.config.EvictedFunc(k, item.value)
		}
	}

	c.init()
}

func (c *simpleCache) Keys() []cache.Key {
	var keys []cache.Key
	for k := range c.cache {
		if _, ok := c.lookup(k); !ok {
			continue
		}

		keys = append(keys, k)
	}
	return keys
}

func (c *simpleCache) Len() int {
	return len(c.GetALL())
}
