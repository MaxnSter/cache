package lru

import (
	"github.com/MaxnSter/cache"
	"time"
)

const cacheName = "lru"

type lruCache struct {
}

func init() {
	cache.RegisterCreator(cacheName, newLruCache)
}

func newLruCache(config *cache.Config) cache.Cache {
	panic("implement me")
	return &lruCache{}
}

func (c *lruCache) Set(cache.Key, cache.Value) error {
	panic("implement me")
}

func (c *lruCache) SetWithExpire(cache.Key, cache.Value, time.Duration) error {
	panic("implement me")
}

func (c *lruCache) Get(cache.Key) (cache.Value, bool) {
	panic("implement me")
}

func (c *lruCache) GetALL() map[cache.Key]cache.Value {
	panic("implement me")
}

func (c *lruCache) Load(cache.Key) error {
	panic("implement me")
}

func (c *lruCache) LoadWithExpire(cache.Key, time.Duration) error {
	panic("implement me")
}

func (c *lruCache) Remove(cache.Key) bool {
	panic("implement me")
}

func (c *lruCache) Purge() {
	panic("implement me")
}

func (c *lruCache) Keys() []cache.Key {
	panic("implement me")
}

func (c *lruCache) Len() int {
	panic("implement me")
}

func (c *lruCache) HitCount() uint64 {
	panic("implement me")
}

func (c *lruCache) MissCount() uint64 {
	panic("implement me")
}

func (c *lruCache) LookupCount() uint64 {
	panic("implement me")
}

func (c *lruCache) HitRate() float64 {
	panic("implement me")
}
