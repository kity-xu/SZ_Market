package models

import (
	"haina.com/share/logging"
	"haina.com/share/models"

	"haina.com/market/f9/config"
)

func init() {
	cfg := config.Default(APP_PID)

	//初始化mysql配置
	logging.Info("初始化mysql_1配置")
	logging.Info(cfg.Db.DriveName)
	logging.Info(cfg.Db.DataSource)

	logging.Info("初始化mysql_2配置")
	logging.Info(cfg.Db1.DriveName)

	logging.Info(cfg.Db1.DataSource)
	err := models.Init(cfg.Db.DriveName, cfg.Db.DataSource)
	err = models.InitDB(cfg.Db1.DriveName, cfg.Db1.DataSource)
	if err != nil {
		logging.Fatal(err)
	}
	//初始化redis配置
	logging.Info("初始化redis_1配置")
	logging.Info("初始化redis_2配置")
	logging.Info(cfg.Redis.Addr)
	//	redis.Init(
	//		cfg.Redis.Addr,
	//		cfg.Redis.Db,
	//		cfg.Redis.Auth,
	//		cfg.Redis.Timeout,
	//	)
	InitRedisFrame(&cfg.Redis, &cfg.RedisCache)

	//	logging.Info(cfg.RedisCache.Addr)
	//	//	redis.Init(
	//	//		cfg.Redis.Addr,
	//	//		cfg.Redis.Db,
	//	//		cfg.Redis.Auth,
	//	//		cfg.Redis.Timeout,
	//	//	)
	//	InitRedisFrame(&cfg.RedisCache, &cfg.RedisCache)

	//RedisCache = redis.NewRedisPool(cfg.Redis.Addr, cfg.Redis.Db, cfg.Redis.Auth, cfg.Redis.Timeout)

	//	RedisCache = NewRedisPool(cfg.Redis.Addr, cfg.Redis.Db, cfg.Redis.Auth, cfg.Redis.Timeout)
	//	RedisStore = NewRedisPool(data_source)
	//	RedisML = NewRedisPool(microlink)
}
