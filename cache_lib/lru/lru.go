package lru

import (
	"container/list"
	"fmt"
	"github.com/MaxnSter/cache"
	"github.com/jonboulle/clockwork"
	"golang.org/x/sync/singleflight"
	"sync"
	"time"
)

const cacheName = "lru"

type cacheItem struct {
	key        cache.Key
	value      cache.Value
	clock      clockwork.Clock
	expiration *time.Time
}

type lruCache struct {
	*cache.Stats
	l     *sync.Mutex
	group *singleflight.Group

	config cache.Config
	list   *list.List
	cache  map[cache.Key]*list.Element
	clock  clockwork.Clock
}

func init() {
	cache.RegisterCreator(cacheName, newLruCache)
}

func newLruCache(config *cache.Config) cache.Cache {
	c := &lruCache{
		Stats:  &cache.Stats{},
		l:      &sync.Mutex{},
		config: *config,
		list:   list.New(),
		clock:  clockwork.NewRealClock(),
		group:  &singleflight.Group{},
	}

	c.init()
	return c
}

func (c lruCache) init() {
	if c.config.Size > 0 {
		c.cache = make(map[cache.Key]*list.Element, c.config.Size)
		return
	}

	c.cache = map[cache.Key]*list.Element{}
}

func (c *lruCache) Set(k cache.Key, v cache.Value) error {
	c.l.Lock()
	defer c.l.Unlock()

	return c.setLocked(k, v)
}

func (c *lruCache) setLocked(k cache.Key, v cache.Value) error {
	var err error
	if c.config.Serializer != nil {
		v, err = c.config.Serializer.Serialize(k, v)
		if err != nil {
			return err
		}
	}

	element, ok := c.cache[k]
	if ok {
		element.Value.(*cacheItem).value = v
		c.list.MoveToFront(element)
	} else {
		if c.config.Size > 0 && len(c.cache) >= c.config.Size {
			c.evictLocked(1)
		}

		element = &list.Element{
			Value: &cacheItem{
				key:   k,
				value: v,
				clock: c.clock,
			},
		}

		c.list.PushFront(element)
		c.cache[k] = element
	}

	if c.config.Expiration != nil {
		t := c.clock.Now().Add(*c.config.Expiration)
		element.Value.(*cacheItem).expiration = &t
	}

	if c.config.AddedFunc != nil {
		c.config.AddedFunc(k, v)
	}

	return nil
}

func (c *lruCache) SetWithExpire(k cache.Key, v cache.Value, t time.Duration) error {
	c.l.Lock()
	defer c.l.Unlock()

	if err := c.setLocked(k, v); err != nil {
		return err
	}

	nt := c.clock.Now().Add(t)
	c.cache[k].Value.(*cacheItem).expiration = &nt
	return nil
}

func (c *lruCache) Get(k cache.Key) (cache.Value, bool) {
	item, ok := c.lookup(k)
	if !ok {
		c.IncrMissCount()
		return nil, false
	}

	var err error
	var v cache.Value
	if c.config.Deserializer != nil {
		v, err = c.config.Deserializer.Deserialize(k, item.value)
		if err != nil {
			return nil, false
		}
	}
	c.IncrHitCount()
	return v, true
}

func (c *lruCache) lookup(k cache.Key) (*cacheItem, bool) {
	c.l.Lock()
	defer c.l.Unlock()

	element, ok := c.cache[k]
	if !ok {
		return nil, false
	}

	item := element.Value.(*cacheItem)
	if item.expiration != nil && c.clock.Now().After(*item.expiration) {
		c.removeLocked(k)
		return nil, false
	}

	return item, true
}

func (c *lruCache) GetALL() map[cache.Key]cache.Value {
	m := map[cache.Key]cache.Value{}
	for k := range c.cache {
		item, ok := c.lookup(k)
		if !ok {
			continue
		}

		m[k] = item.value
	}
	return m
}

func (c *lruCache) Load(k cache.Key) error {
	return c.load(k, func(key cache.Key, value cache.Value) error {
		return c.Set(key, value)
	})
}

func (c *lruCache) LoadWithExpire(k cache.Key, t time.Duration) error {
	return c.load(k, func(key cache.Key, value cache.Value) error {
		return c.SetWithExpire(key, value, t)
	})
}

func (c *lruCache) load(k cache.Key, loaded func(cache.Key, cache.Value) error) error {
	_, err, _ := c.group.Do(k, func() (interface{}, error) {
		if item, ok := c.lookup(k); ok {
			return item.value, nil
		}

		if c.config.Loader == nil {
			return nil, fmt.Errorf("loader not found")
		}

		v, err := c.config.Loader.Load(k)
		if err != nil {
			return nil, err
		}

		if err := loaded(k, v); err != nil {
			return nil, err
		}

		return v, nil
	})

	return err
}

func (c *lruCache) Remove(k cache.Key) bool {
	c.l.Lock()
	defer c.l.Unlock()

	return c.removeLocked(k)
}

func (c *lruCache) evictLocked(count int) {
	var evicted int
	t := c.clock.Now()

	for k, element := range c.cache {
		if evicted >= count {
			return
		}

		item := element.Value.(*cacheItem)
		if item.expiration != nil && t.After(*item.expiration) {
			c.removeLocked(k)
			evicted++
		}
	}

	for i := evicted; i < count; i++ {
		c.removeLocked(c.list.Back().Value.(cacheItem).key)
	}
}

func (c *lruCache) removeLocked(k cache.Key) bool {
	element, ok := c.cache[k]
	if !ok {
		return false
	}

	delete(c.cache, k)
	c.list.Remove(element)
	if c.config.EvictedFunc != nil {
		c.config.EvictedFunc(k, element.Value)
	}
	return true
}

func (c *lruCache) Purge() {
	panic("implement me")
}

func (c *lruCache) Keys() []cache.Key {
	var keys []cache.Key
	for k := range c.cache {
		if _, ok := c.lookup(k); !ok {
			continue
		}

		keys = append(keys, k)
	}
	return keys
}

func (c *lruCache) Len() int {
	return len(c.GetALL())
}
