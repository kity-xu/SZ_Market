package models

import "errors"

// App Setting
//---------------------------------------------------------------------------------
const (
	APP_NAME    = "hqpost"
	APP_VERSION = "0.1.0.4"
	APP_PID     = "hqpost"
)

// global_table
//---------------------------------------------------------------------------------
const (
	REDISKEY_SECURITY_NAME_ID = "hq:st:name:%d" ///<证券代码(参数：sid) (hq-init写入)
)

var (
	ERROR_REDIS_LIST_NULL = errors.New("redis list is null")
)
