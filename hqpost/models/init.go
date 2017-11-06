package models

import (
	"haina.com/market/hqpost/config"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

var (
	RedisStore *redis.RedisPool // 数据源   存放着从数据源生产者产生的数据(pb,bin)
)

func NewRedisPool(r *config.RedisStore) *redis.RedisPool {
	return redis.NewRedisPool(r.Addr, r.Db, r.Auth, r.Timeout)
}

// 初始化 Redis 配置(应答缓存Cache + 数据Store 架构)
//   根据架构设计进行灵活调配
func InitRedisFrame(data_source *config.RedisStore) {
	if data_source == nil {
		logging.Fatal(" InitRedisFrame failed !!!")
	}
	RedisStore = NewRedisPool(data_source)
}

func init() {

	cfg := config.Default(APP_PID)
	// 初始化 Redis 配置
	redis.Init(
		cfg.Redis.Addr,
		cfg.Redis.Db,
		cfg.Redis.Auth,
		cfg.Redis.Timeout)

	InitRedisFrame(&cfg.Redis)
}
