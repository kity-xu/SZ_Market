package models

import (
	"haina.com/market/hqpublish/config"
)

var (
	FStore *config.FileStore
	TTL    *config.CacheTTL
)

func init() {

	cfg := config.Default(APP_PID)
	FStore = &cfg.File
	TTL = &cfg.TTL

	// 初始化 Redis 配置
	InitRedisFrame(&cfg.RedisCache, &cfg.Redis)
}
