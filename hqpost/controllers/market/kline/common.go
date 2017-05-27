package kline

import (
	"ProtocolBuffer/format/kline"

	"haina.com/market/hqpost/config"
)

var cfg *config.AppConfig

/*****************************************Const***************************************/
//Exchange 交易所
const (
	SH = "/sh/"
	SZ = "/sz/"
)

//// K线历史部分
const (
	REDISKEY_SECURITY_MIN    = "hq:st:min:%d"    ///<证券分钟线数据(参数：sid) (calc写入)
	REDISKEY_SECURITY_HDAY   = "hq:st:hday:%d"   ///<证券历史日K线数据(参数：sid) (hq-post写入)
	REDISKEY_SECURITY_HWEEK  = "hq:st:hweek:%d"  ///<证券周K线(参数：sid)
	REDISKEY_SECURITY_HMONTH = "hq:st:hmonth:%d" ///<证券月K线(参数：sid)
	REDISKEY_SECURITY_HYEAR  = "hq:st:hyear:%d"  ///<证券年K线(参数：sid)
)

//K线、指数定义
type Security struct {
	sids *[]int32
	list SecurityList
}

//单个股票
type SingleSecurity struct {
	Sid int32 //股票SID

	Date      []int32               //单个股票的历史日期
	SigStock  map[int32]kline.KInfo //单个股票的历史数据
	WeekDays  *[][]int32            //单个股票的历史周天
	MonthDays *[][]int32            //单个股票的历史月天
	YearDays  *[][]int32            //单个股票的历史年天
	today     *kline.KInfo          //单个股票的当天数据
}

//所有股票
type SecurityList struct {
	Securitys *[]SingleSecurity
}
