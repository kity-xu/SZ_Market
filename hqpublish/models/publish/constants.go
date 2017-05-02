package publish

import (
	"errors"
)

var (
	ERROR_REDIS_LIST_NULL  = errors.New("redis list is null")
	ERROR_KLINE_BEGIN_TIME = errors.New("error's kline begin time")
	ERROR_KLINE_DATA_NULL  = errors.New("kline data is null")
)

/// 证券代码表部分
const ( // 上游定义已修改 作废
	REDISKEY_SECURITY_CODETABLE = "hgs:global:securitytable"   ///<证券代码表（全）
	REDISKEY_SECURITY_INFO_ID   = "hgs:global:securityinfo:%d" ///<证券代码(参数：sid)
	REDISKEY_SECURITY_INFO_CODE = "hgs:global:securityinfo:%s" ///<证券代码(参数：scode)
)
const (
	REDISKEY_MARKET_SECURITY_TABLE = "hq:market:sts:%d" ///<证券市场代码表(参数：MarketID)  (hq-init写入)
	REDISKEY_SECURITY_NAME_ID      = "hq:st:name:%d"    ///<证券代码(参数：sid) (hq-init写入)
	REDISKEY_SECURITY_NAME_CODE    = "hq:st:name:%s"    ///<证券代码(参数：scode) (hq-init写入)
)

/// 分钟线
const (
	REDISKEY_SECURITY_MIN  = "hq:st:min:%d"  ///<证券分钟线数据(参数：sid) (calc写入)
	REDISKEY_SECURITY_HMIN = "hq:st:hmin:%d" ///<证券历史分钟线数据(参数：sid) (hq-post写入)
)

const (
	REDISKEY_SECURITY_CODETABLE_REPLY_TTL = 300
)
