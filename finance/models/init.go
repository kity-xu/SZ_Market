package models

import (
	"haina.com/market/finance/config"
	"haina.com/share/logging"
	"haina.com/share/models"
	"haina.com/share/store/redis"
)

func init() {

	cfg := config.Default(APP_PID)

	//初始化 MySQL 配置
	err := models.Init(cfg.Db.DriverName, cfg.Db.DataSource)
	if err != nil {
		logging.Fatal(err)
		return
	}

	// 初始化 Redis 配置
	redis.Init(
		cfg.Redis.Addr,
		cfg.Redis.Db,
		cfg.Redis.Auth,
		cfg.Redis.Timeout)
}
