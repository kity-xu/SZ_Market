package models

import (
	"strings"

	"haina.com/market/hqpublish/config"
	"haina.com/share/logging"
	"haina.com/share/models"
	"haina.com/share/store/mongo"
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

	// 初始化 MongDB 配置
	i := strings.LastIndexByte(cfg.Mongo.Source, '/')
	if i <= 0 {
		logging.Fatalf("xml mongoStore:source error: %q", cfg.Mongo.Source)
	}
	mdb := cfg.Mongo.Source[i+1:]
	if err := mongo.Init(cfg.Mongo.Source, mdb); err != nil {
		logging.Fatal(err)
	}
}
