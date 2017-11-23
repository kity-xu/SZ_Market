package models

import "errors"

// App Setting
//---------------------------------------------------------------------------------
const (
	APP_NAME    = "hqpost"
	APP_VERSION = "2.0.0"
	APP_PID     = "hqpost"
)

// global_table
//---------------------------------------------------------------------------------
const (
	REDISKEY_SECURITY_NAME_ID = "hq:st:name:%d" ///<证券代码(参数：sid) (hq-init写入)
	REDISKEY_BLOCKINDEX_MIN   = "hq:st:min:%s"  ///板块指数快照
)

var (
	ERROR_REDIS_LIST_NULL = errors.New("redis list is null")
)
