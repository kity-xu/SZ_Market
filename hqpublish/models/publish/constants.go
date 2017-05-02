package publish

import (
	"errors"
)

var (
	REDIS_ERROR_LIST_NULL = errors.New("redis list is null")
)

/// 证券代码表部分
const (
	REDISKEY_SECURITY_CODETABLE = "hgs:global:securitytable"   ///<证券代码表（全）
	REDISKEY_SECURITY_INFO_ID   = "hgs:global:securityinfo:%d" ///<证券代码(参数：sid)
	REDISKEY_SECURITY_INFO_CODE = "hgs:global:securityinfo:%s" ///<证券代码(参数：scode)
)

/// 分钟线
const (
	REDISKEY_SECURITY_MIN  = "hq:st:min:%d"  ///<证券分钟线数据(参数：sid) (calc写入)
	REDISKEY_SECURITY_HMIN = "hq:st:hmin:%d" ///<证券历史分钟线数据(参数：sid) (hq-post写入)
)

const (
	REDISKEY_SECURITY_CODETABLE_REPLY_TTL = 300
)
