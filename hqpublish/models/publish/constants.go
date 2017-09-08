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
	FILE_HMINDATA_NULL   = errors.New("the file is empty")
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

//排序&板块
const (
	REDISKEY_SORT_KDAY_H  = "hq:sort:%d:%d"      ///<排序结果(顺序，参数：证券组id,字段id) (calc写入)
	REDIS_KEY_CACHE_BLOCK = "hq:init:bk:%d"      ///板块
	REDISKEY_ELEMENT      = "hq:init:bk:1100:%d" ///成份股
)

//资金流向
const (
	REDISKEY_FUND_FLOW = "hq:trade:min:%d" ///资金流向
)

// 证券集合(板块)
const (
	REDISKEY_STOCK_BLOCK      = "hq:bk:%d:*"          ///板块信息(参数：板块基础id)(calc写入)
	REDISKEY_STOCK_BLOCK_BASE = "hq:sort:11:%d"       ///板块信息(参数：排序id)(calc写入)
	REDISKEY_STOCK_BLOCK_INFO = "hq:init:bk:%d:*"     ///板块信息(calc写入)
	REDISKEY_STOCK_BLOCK_SHOT = "hq:bk:snap:%d"       ///板块快照信息(参数：板块基础id, 板块id)(calc写入)
	REDISKEY_STOCK_BLOCK_SID  = "hq:init:bk:stock:%d" ///成分股所属板块
)

// 资金统计
const (
	REDISKEY_TRADE_PRICE = "hq:trade:price:%d" ///证券分价统计（参数：sid）（calc写入）
)

// RedisCache: 缓存首页
const (
	REDISKEY_L2CACHE_INDEX_MOBILE = "hq:index:mobile" // 移动端 /api/hq/mindex
	REDISKEY_L2CACHE_INDEX_PC     = "hq:index:pc"     // PC端 /api/hq/pcindex
)

// 基准日 20100101
const (
	Baseday = 20100104 //基准日后第一个工作日
	Workday = 251      //一年中交易日（除去双休日和法定假日---只多不少）
)
