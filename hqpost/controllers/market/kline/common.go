package kline

import (
	"fmt"
	"strings"
	"time"

	"haina.com/share/logging"
)

/*****************************************Const***************************************/
//Exchange 交易所
const (
	SH = "/sh/"
	SZ = "/sz/"
)

//// K线历史部分
const (
	REDISKEY_SECURITY_HDAY   = "hq:st:hday:%d"   ///<证券历史日K线数据(参数：sid) (hq-post写入)
	REDISKEY_SECURITY_HWEEK  = "hq:st:hweek:%d"  ///<证券周K线(参数：sid)
	REDISKEY_SECURITY_HMONTH = "hq:st:hmonth:%d" ///<证券月K线(参数：sid)
	REDISKEY_SECURITY_HYEAR  = "hq:st:hyear:%d"  ///<证券年K线(参数：sid)
	REDISKEY_SECURITY_HMIN   = "hq:st:hmin:%d"   ///<<证券历史分钟线数据(参数：sid) (hq-post写入)
	REDISKEY_SECURITY_HMIN5  = "hq:st:hmin5:%d"  ///<证券5分钟K线(参数：sid)
	REDISKEY_SECURITY_HMIN15 = "hq:st:hmin15:%d" ///<证券15分钟K线(参数：sid)
	REDISKEY_SECURITY_HMIN30 = "hq:st:hmin30:%d" ///<证券30分钟K线(参数：sid)
	REDISKEY_SECURITY_HMIN60 = "hq:st:hmin60:%d" ///<证券60分钟K线(参数：sid)
)

/*****************************************Structs***************************************/

//K线定义
type StockSingle struct {
	SID    int32  // 证券ID
	Time   int32  // 时间 unix time
	PreCPx int32  // 昨收价 * 10000
	OpenPx int32  // 开盘价 * 10000
	HighPx int32  // 最高价 * 10000
	LowPx  int32  // 最低价 * 10000
	LastPx int32  // 最新价 * 10000
	Volume int64  // 成交量
	Value  int64  // 成交额 * 10000
	AvgPx  uint32 // 平均价 * 10000

}

type Security struct {
	list SecurityList
}

//单个股票
type SingleSecurity struct {
	Sid       int32                 //股票SID
	Date      []int32               //单个股票的历史日期
	SigStock  map[int32]StockSingle //单个股票的历史数据
	WeekDays  *[][]int32            //单个股票的周天
	MonthDays *[][]int32            //单个股票的月天
	YearDays  *[][]int32            //单个股票的年天
}

//所有股票
type SecurityList struct {
	Securitys *[]SingleSecurity
}

/*****************************************Functions***************************************/
func NewSecurity() *Security {
	return &Security{}
}

func IntToTime(date int) time.Time {
	swap := date % 10000
	year := date / 10000
	month := swap / 100
	day := swap % 100
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func DateAdd(date int) time.Time {
	var sat time.Time
	swap := date % 10000
	year := date / 10000
	month := swap / 100
	day := swap % 100

	//logging.Debug("%d-%d-%d", year, month, day)

	baseTime := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	weekday := baseTime.Weekday().String()

	var basedate string
	if strings.EqualFold(weekday, "Monday") {
		basedate = fmt.Sprintf("%d%s", 24*5, "h")

	} else if strings.EqualFold(weekday, "Tuesday") {
		basedate = fmt.Sprintf("%d%s", 24*4, "h")

	} else if strings.EqualFold(weekday, "Wednesday") {
		basedate = fmt.Sprintf("%d%s", 24*3, "h")

	} else if strings.EqualFold(weekday, "Thursday") {
		basedate = fmt.Sprintf("%d%s", 24*2, "h")

	} else if strings.EqualFold(weekday, "Friday") {
		basedate = fmt.Sprintf("%d%s", 24*1, "h")

	} else {
		logging.Error("Invalid trade date...")
		return sat
	}

	dd, _ := time.ParseDuration(basedate)
	sat = baseTime.Add(dd) //Saturday（星期六）
	return sat
}

func IntToMonth(date int) time.Time {
	swap := date % 10000
	year := date / 10000
	month := swap / 100

	return time.Date(year, time.Month(month), 0, 0, 0, 0, 0, time.UTC)
}
