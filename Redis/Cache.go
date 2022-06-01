package Redis

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var GlobalCache *cache.Cache

func CacheInit() {
	cache1 := cache.New(60*time.Second, 10*time.Second)
	GlobalCache = cache1
}
