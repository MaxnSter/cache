package lfu

import (
	"github.com/MaxnSter/cache"
	"time"
)

const cacheName = "lfu"

type lfuCache struct {
}

func init() {
	cache.RegisterCreator(cacheName, newLfuCache)
}

func newLfuCache(config *cache.Config) cache.Cache {
	panic("implement me")
	return &lfuCache{}
}

func (c *lfuCache) Set(cache.Key, cache.Value) error {
	panic("implement me")
}

func (c *lfuCache) SetWithExpire(cache.Key, cache.Value, time.Duration) error {
	panic("implement me")
}

func (c *lfuCache) Get(cache.Key) (cache.Value, bool) {
	panic("implement me")
}

func (c *lfuCache) GetALL() map[cache.Key]cache.Value {
	panic("implement me")
}

func (c *lfuCache) Load(cache.Key) error {
	panic("implement me")
}

func (c *lfuCache) LoadWithExpire(cache.Key, time.Duration) error {
	panic("implement me")
}

func (c *lfuCache) Remove(cache.Key) bool {
	panic("implement me")
}

func (c *lfuCache) Purge() {
	panic("implement me")
}

func (c *lfuCache) Keys() []cache.Key {
	panic("implement me")
}

func (c *lfuCache) Len() int {
	panic("implement me")
}

func (c *lfuCache) HitCount() uint64 {
	panic("implement me")
}

func (c *lfuCache) MissCount() uint64 {
	panic("implement me")
}

func (c *lfuCache) LookupCount() uint64 {
	panic("implement me")
}

func (c *lfuCache) HitRate() float64 {
	panic("implement me")
}
