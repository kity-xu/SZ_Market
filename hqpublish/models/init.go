package models

import (
	"haina.com/market/hqpublish/config"
	"haina.com/share/store/redis"
)

var (
	FStore *config.FileStore
	TTL    *config.CacheTTL
)

func init() {

	cfg := config.Default(APP_PID)
	FStore = &cfg.File
	TTL = &cfg.TTL

	// 初始化 Redis 配置（即将废弃）
	redis.Init(
		cfg.Redis.Addr,
		cfg.Redis.Db,
		cfg.Redis.Auth,
		cfg.Redis.Timeout)

	InitRedisFrame(&cfg.RedisCache, &cfg.Redis)
}
