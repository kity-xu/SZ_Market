package models

import (
	"haina.com/market/hqpublish/config"
	"haina.com/share/store/redis"
)

func init() {

	cfg := config.Default(APP_PID)

	// 初始化 Redis 配置（即将废弃）
	redis.Init(
		cfg.Redis.Addr,
		cfg.Redis.Db,
		cfg.Redis.Auth,
		cfg.Redis.Timeout)

	InitRedisFrame(&cfg.RedisCache, &cfg.Redis)
}
