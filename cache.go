package cache

import (
	"time"
)

type (
	Key   = string
	Value = interface{}
)

type Cache interface {
	Set(Key, Value) error
	SetWithExpire(Key, Value, time.Duration) error

	Get(Key) (Value, bool)
	GetALL() map[Key]Value

	Load(Key) error
	LoadWithExpire(Key, time.Duration) error

	Remove(Key) bool
	Purge()

	Keys() []Key
	Len() int

	statsAccessor
}

type Loader interface {
	Load(Key) (Value, error)
}

type Serializer interface {
	Serialize(Key, Value) (Value, error)
}

type Deserializer interface {
	Deserialize(Key, Value) (Value, error)
}

type LoaderAdaptor func(Key) (Value, error)

func (a LoaderAdaptor) Load(k Key) (Value, error) {
	return a(k)
}

type SerializerAdaptor func(Key, Value) (Value, error)

func (a SerializerAdaptor) Serialize(k Key, v Value) (Value, error) {
	return a(k, v)
}

type DeserializeAdaptor func(Key, Value) (Value, error)

func (a DeserializeAdaptor) Deserialize(k Key, v Value) (Value, error) {
	return a(k, v)
}

type (
	LoaderFunc       func(Key) (Value, error)
	EvictedFunc      func(Key, Value)
	PurgeVisitorFunc func(Key, Value)
	AddedFunc        func(Key, Value)
	DeserializeFunc  func(Key, Value) (Value, error)
	SerializeFunc    func(Key, Value) (Value, error)
)

type Config struct {
	Size             int
	EvictedFunc      EvictedFunc
	PurgeVisitorFunc PurgeVisitorFunc
	AddedFunc        AddedFunc
	Expiration       *time.Duration
	Loader           Loader
	Serializer       Serializer
	Deserializer     Deserializer
}

func NewCache(cacheType string, config *Config) Cache {
	return MustGetCacheCreator(cacheType)(config)
}

type OptionFunc func(*Config)

func NewCacheConfig(options ...OptionFunc) *Config {
	c := &Config{}

	for _, o := range options {
		o(c)
	}
	return c
}

func WithSize(size int) func(*Config) {
	return func(config *Config) {
		config.Size = size
	}
}

func WithExpiration(t time.Duration) func(*Config) {
	return func(config *Config) {
		config.Expiration = &t
	}
}

func WithEvictedFunc(f EvictedFunc) func(*Config) {
	return func(config *Config) {
		config.EvictedFunc = f
	}
}

func WithPrugeVisitorFunc(f PurgeVisitorFunc) func(*Config) {
	return func(config *Config) {
		config.PurgeVisitorFunc = f
	}
}

func WithAddedFunc(f AddedFunc) func(*Config) {
	return func(config *Config) {
		config.AddedFunc = f
	}
}

func WithLoader(l Loader) func(*Config) {
	return func(config *Config) {
		config.Loader = l
	}
}

func WithSerializer(s Serializer) func(*Config) {
	return func(config *Config) {
		config.Serializer = s
	}
}

func WithDeserializer(s Deserializer) func(*Config) {
	return func(config *Config) {
		config.Deserializer = s
	}
}
