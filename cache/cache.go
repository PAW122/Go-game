package cache

import (
	"strings"

	"github.com/patrickmn/go-cache"
)

var c *cache.Cache

func init() {
	c = cache.New(cache.NoExpiration, cache.NoExpiration)
}

// read
func GetCache(key string) (interface{}, bool) {
	cachedData, found := c.Get(key)
	return cachedData, found
}

// getAllCacheData returns all keys and their corresponding data stored in the cache
func GetAllCacheData() map[string]interface{} {
	allData := make(map[string]interface{})

	// Iterate over all items in the cache
	for k, v := range c.Items() {
		allData[k] = v.Object
	}

	return allData
}

func GetCacheSize() int64 {
	res := int64(c.ItemCount())
	return res
}

func SaveCache(key string, data interface{}) bool {
	duration := cache.NoExpiration
	c.Set(key, data, duration)

	// Unieważnij wszystkie bardziej ogólne klucze
	invalidateParentKeys(key)

	return true
}

func DeleteCache(key string) bool {
	c.Delete(key)
	return true
}

func invalidateParentKeys(key string) {
	parts := strings.Split(key, ".")
	for i := 1; i < len(parts); i++ {
		parentKey := strings.Join(parts[:i], ".")
		c.Delete(parentKey)
	}
}
