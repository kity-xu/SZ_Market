package models

import (
	"fmt"

	"haina.com/market/hqpublish/config"
	"haina.com/share/store/redis"
)

func init() {

	cfg := config.Default(APP_PID)

	// 初始化 Redis 配置
	redis.Init(
		cfg.Redis.Addr,
		cfg.Redis.Db,
		cfg.Redis.Auth,
		cfg.Redis.Timeout)
	fmt.Printf("%v\n", cfg.RedisCache)
}
