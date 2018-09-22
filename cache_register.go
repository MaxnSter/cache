package cache

import "fmt"

type Creator func(*Config) Cache

var creators = map[string]Creator{}

func RegisterCreator(cacheType string, creator Creator) {
	if _, ok := creators[cacheType]; ok {
		panic(fmt.Sprintf("dup register cache type:%s", cacheType))
	}

	creators[cacheType] = creator
}

func MustGetCacheCreator(cacheType string) Creator {
	if _, ok := creators[cacheType]; !ok {
		panic(fmt.Sprintf("cache type :%s, dit not register", cacheType))
	}

	return creators[cacheType]
}
