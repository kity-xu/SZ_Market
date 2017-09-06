package controllers

import (
	"errors"
)

const (
	/// 证券代码表
	REDISKEY_MARKET_SECURITY_TABLE = "hq:market:sts:%d"      ///<证券市场代码表(参数：MarketID)  (hq-init写入)
	REDISKEY_SECURITY_NAME_ID      = "hq:st:name:%d"         ///<证券代码(参数：sid) (hq-init写入)
	REDISKEY_SECURITY_NAME_CODE    = "hq:st:name:%s"         ///<证券代码(参数：scode) (hq-init写入)
	REDISKEY_HQPOST_EXECUTED_TIME  = "hq:post:time:executed" ///hqpost 上次执行完毕时的时间戳

	REDISKEY_SECURITY_SNAP = "hq:st:snap:%d" ///<证券快照数据(参数：sid) (calc写入)
)

var (
	ERROR_INDEX_MAYBE_OUTOF_RANGE = errors.New("index maybe out of range")
)

// 字符串型数据长度定义

const (
	SECURITY_CODE_LEN = 24 ///< 证券代码长度
	SECURITY_NAME_LEN = 40 ///< 证券名称长度
	SECURITY_DESC_LEN = 8  ///< 英文简称
	INDUSTRY_CODE_LEN = 8  ///< 行业代码
	SECURITY_ISIN_LEN = 16 ///< 证券国际代码信息

)
