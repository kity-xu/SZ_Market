package publish

import (
	"errors"
)

var (
	ERROR_REQUEST_PARAM    = errors.New("request param error")
	ERROR_REDIS_DATE_NULL  = errors.New("redis date is null")
	ERROR_REDIS_LIST_NULL  = errors.New("redis list is null")
	ERROR_KLINE_BEGIN_TIME = errors.New("error's kline begin time")
	ERROR_KLINE_DATA_NULL  = errors.New("kline data is null")
	ERROR_INVALID_DATA     = errors.New("Invalid data or data is null")

	INVALID_FILE_PATH    = errors.New("Invalid file path")
	INVALID_REQUEST_PARA = errors.New("Invalid request parameter type")

	READ_REDIS_STORE_NULL = errors.New("redis store is null")
)

/// 证券代码表部分
const ( // 上游定义已修改 作废
	REDISKEY_SECURITY_CODETABLE = "hgs:global:securitytable"   ///<证券代码表（全）
	REDISKEY_SECURITY_INFO_ID   = "hgs:global:securityinfo:%d" ///<证券代码(参数：sid)
	REDISKEY_SECURITY_INFO_CODE = "hgs:global:securityinfo:%s" ///<证券代码(参数：scode)
)

const (
	REDISKEY_MARKET_STATUS = "hq:market:%d" ///<证券市场状态(参数：MarketID) (calc写入)
)

const (
	REDISKEY_MARKET_SECURITY_TABLE_ASTOCK = "hq:market:sts:%s" ///A股市场
	REDISKEY_MARKET_SECURITY_TABLE        = "hq:market:sts:%d" ///<证券市场代码表(参数：MarketID)  (hq-init写入)
	REDISKEY_SECURITY_NAME_ID             = "hq:st:name:%d"    ///<证券代码(参数：sid) (hq-init写入)
	REDISKEY_SECURITY_NAME_CODE           = "hq:st:name:%s"    ///<证券代码(参数：scode) (hq-init写入)
	REDISKEY_SECURITY_STATIC              = "hq:st:static:%d"  ///<证券静态数据(参数：sid) (hq-init写入)
)

/// 证券快照
const (
	REDISKEY_SECURITY_SNAP = "hq:st:snap:%d" ///<证券快照数据(参数：sid) (calc写入)
)

/// 分钟线
const (
	REDISKEY_SECURITY_MIN = "hq:st:min:%d" ///<证券分钟线数据(参数：sid) (calc写入)
)

/// 分笔成交统计
const (
	REDISKEY_SECURITY_TRADE = "hq:st:trade:%d" ///<证券分笔成交数据(参数：sid) (calc写入)
)

//历史K线
const (
	REDISKEY_SECURITY_HDAY   = "hq:st:hday:%d"   ///<证券历史日K线数据(参数：sid) (hq-post写入)
	REDISKEY_SECURITY_HWEEK  = "hq:st:hweek:%d"  ///<证券周K线(参数：sid)
	REDISKEY_SECURITY_HMONTH = "hq:st:hmonth:%d" ///<证券月K线(参数：sid)
	REDISKEY_SECURITY_HYEAR  = "hq:st:hyear:%d"  ///<证券年K线(参数：sid)
	REDISKEY_SECURITY_HMIN   = "hq:st:hmin:%d"   ///<<证券历史分钟线数据(参数：sid) (hq-post写入)
	REDISKEY_SECURITY_HMIN1  = "hq:st:hmin1:%d"  ///<证券1分钟K线(hq-post写入)
	REDISKEY_SECURITY_HMIN5  = "hq:st:hmin5:%d"  ///<证券5分钟K线(参数：sid)
	REDISKEY_SECURITY_HMIN15 = "hq:st:hmin15:%d" ///<证券15分钟K线(参数：sid)
	REDISKEY_SECURITY_HMIN30 = "hq:st:hmin30:%d" ///<证券30分钟K线(参数：sid)
	REDISKEY_SECURITY_HMIN60 = "hq:st:hmin60:%d" ///<证券60分钟K线(参数：sid)
)

//排序
const (
	REDISKEY_SORT_KDAY_H = "hq:sort:%d:%d"    ///<排序结果(顺序，参数：证券组id,字段id) (calc写入)
	REDISKEY_BLOCK       = "hq:init:bk:%d"    ///板块
	REDISKEY_ELEMENT     = "hq:init:bk:%d:%d" ///成份股
)
