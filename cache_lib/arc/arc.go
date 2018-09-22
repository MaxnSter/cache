package arc

import (
	"github.com/MaxnSter/cache"
	"time"
)

const cacheName = "arc"

type arcCache struct {
}

func init() {
	cache.RegisterCreator(cacheName, newArcCache)
}

func newArcCache(config *cache.Config) cache.Cache {
	panic("implement me")
	return &arcCache{}
}

func (c *arcCache) Set(cache.Key, cache.Value) error {
	panic("implement me")
}

func (c *arcCache) SetWithExpire(cache.Key, cache.Value, time.Duration) error {
	panic("implement me")
}

func (c *arcCache) Get(cache.Key) (cache.Value, bool) {
	panic("implement me")
}

func (c *arcCache) GetALL() map[cache.Key]cache.Value {
	panic("implement me")
}

func (c *arcCache) Load(cache.Key) error {
	panic("implement me")
}

func (c *arcCache) LoadWithExpire(cache.Key, time.Duration) error {
	panic("implement me")
}

func (c *arcCache) Remove(cache.Key) bool {
	panic("implement me")
}

func (c *arcCache) Purge() {
	panic("implement me")
}

func (c *arcCache) Keys() []cache.Key {
	panic("implement me")
}

func (c *arcCache) Len() int {
	panic("implement me")
}

func (c *arcCache) HitCount() uint64 {
	panic("implement me")
}

func (c *arcCache) MissCount() uint64 {
	panic("implement me")
}

func (c *arcCache) LookupCount() uint64 {
	panic("implement me")
}

func (c *arcCache) HitRate() float64 {
	panic("implement me")
}
