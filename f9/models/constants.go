package models

import (
	"errors"
)

const (
	APP_NAME    = "market_finance"
	APP_VERSION = "0.0.3.1"
	APP_PID     = "market_finance"
)

const (
	CONTEXT_SCODE = "scode" //证券代码

)

// F10财务分析接口URL请求相关参数
const (
	CONTEXT_TYPE = "type" // 报表类型
	CONTEXT_PAGE = "page" // 当前页码
)

// 数据库表
const (
	TABLE_TQ_OA_STCODE    = "TQ_OA_STCODE"
	TABLE_TQ_SK_BASICINFO = "TQ_SK_BASICINFO"
)

// redis 键值
const (
	REDISKEY_SYSMBOL_BASIC = "finchina:symbol:%d:basic"

	REDISKEY_SNAP          = "hq:st:snap:%v" //快照
	REDISKEY_MARKET_STATUS = "hq:market:%v"  ///<证券市场状态(参数：MarketID) (calc写入)
)

// ERROR
var (
	ERROR_COMPCODE_NULL = errors.New("COMPCODE is NULL")
)
