package models

import (
	"errors"
)

// App Setting
//---------------------------------------------------------------------------------
const (
	APP_NAME    = "market_hqpublish"
	APP_VERSION = "2.0.0"
	APP_PID     = "market_hqpublish"
)

const (
	CONTEXT_FORMAT = "format"       // 数据格式
	ACCESS_TOKEN   = "access_token" // access_token
	CONTEXt_MARKET = "marketID"     // 市场ID
	CONTEXT_SNID   = "snid"         // snid
)

const (
	REDISKEY_HQPOST_EXECUTED_TIME = "hq:post:time:executed" ///hqpost 上次执行完毕时的时间戳
)

//MicroLink 数据库（123.56.30.141）表名
const (
	TABLE_HN_OPT_STOCK    = "hn_opt_stock"    //自选股
	TABLE_HN_OPT_STOCKLOG = "hn_opt_stocklog" //自选股操作日志表
)

//MicroLink Redis 键值
const (
	REDIS_ACCESS_TOKEN_MEMBERID = "m:token:%s" //会员ID
)

var (
	REDIS_MEMBERID_NOT_FIND = errors.New("The member id is not found")
	MYSQL_NOT_FIND          = errors.New("Mysql not found")
)

// 20171212 zxw add ----------------------------begin
const (
	CONTEXT_SCODE = "scode" // 证券代码
)

// F10财务分析接口URL请求相关参数
const (
	CONTEXT_TYPE    = "type"    // 报表类型
	CONTEXT_PERPAGE = "perpage" // 每页条数
	CONTEXT_PAGE    = "page"    // 当前页码
)

// Redis 键值
const (
	REDIS_F10_HOMEPAGE = "hq:f10:homepage:%v" // F10首页数据
	REDIS_F10_COMINFO  = "hq:f10:cominfo:%v"  // 公司详细信息
)

// 20171212 zxw add ----------------------------end
