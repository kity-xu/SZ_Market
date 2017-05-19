package models

import (
	"fmt"

	//grr "github.com/garyburd/redigo/redis"
	"haina.com/market/hqpublish/config"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

var _ = fmt.Print
var _ = logging.Info

var (
	redis_store *redis.RedisPool // 数据源   存放着从数据源生产者产生的数据(pb,bin)
	redis_cache *redis.RedisPool // 应答缓存 缓存PB,json,bin 等直接返回给客户端的内容
)

func NewRedisPool(r *config.RedisStore) *redis.RedisPool {
	return redis.NewRedisPool(r.Addr, r.Db, r.Auth, r.Timeout)
}

// 初始化 Redis 配置(应答缓存Cache + 数据Store 架构)
//   根据架构设计进行灵活调配
func InitRedisFrame(response_cache *config.RedisStore, data_source *config.RedisStore) {
	redis_store = NewRedisPool(data_source)
	redis_cache = NewRedisPool(response_cache)
}

//------------------------------------------------------------------------------
const (
	CACHE_TYPE_SUFFIX_PB   = ".p"
	CACHE_TYPE_SUFFIX_JSON = ".j"
	CACHE_TYPE_SUFFIX_BIN  = ".b"
)

//------------------------------------------------------------------------------
// 从数据Redis获取原始数据
//  数据是从 calc 以及 hqpost, hqinit 写入的不同数据内容区分
//  目前格式有 PB binary
//    其中
//      calc 写入的数据基本都是binary(主要出于效率考虑)
//      hqinit hqpost 写入的数据基本都是PB(出于整存整取和节省空间考虑)
func GetStore(key string) ([]byte, error) {
	return redis_store.GetBytes(key)
}

//------------------------------------------------------------------------------
// 从应答缓存Redis获取键值
// 原始
func GetCache(key string) ([]byte, error) {
	return redis_cache.GetBytes(key)
}

// PB
func GetCachePB(key string) ([]byte, error) {
	key += CACHE_TYPE_SUFFIX_PB
	return GetCache(key)
}

// Json
func GetCacheJson(key string) ([]byte, error) {
	key += CACHE_TYPE_SUFFIX_JSON
	return GetCache(key)
}

// Bin
func GetCacheBin(key string) ([]byte, error) {
	key += CACHE_TYPE_SUFFIX_BIN
	return GetCache(key)
}

//------------------------------------------------------------------------------
func SetCache(key string, ttl int, data []byte) error {
	if ttl < 1 {
		return redis_cache.Setex(key, 3, data)
	} else {
		return redis_cache.Setex(key, ttl, data)
	}
}
func SetCachePB(key string, ttl int, data []byte) error {
	key += CACHE_TYPE_SUFFIX_PB
	return SetCache(key, ttl, data)
}
func SetCacheJson(key string, ttl int, data []byte) error {
	key += CACHE_TYPE_SUFFIX_JSON
	return SetCache(key, ttl, data)
}
func SetCacheBin(key string, ttl int, data []byte) error {
	key += CACHE_TYPE_SUFFIX_BIN
	return SetCache(key, ttl, data)
}
