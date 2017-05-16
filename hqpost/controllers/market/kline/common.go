package kline

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	pbk "ProtocolBuffer/format/kline"
	"errors"

	"haina.com/market/hqpost/config"
	"haina.com/share/logging"
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

/*****************************************Structs***************************************/

//K线、指数定义
type Security struct {
	sids *[]int32
	list SecurityList
}

//单个股票
type SingleSecurity struct {
	Sid       int32               //股票SID
	Date      []int32             //单个股票的历史日期
	SigStock  map[int32]pbk.KInfo //单个股票的历史数据
	WeekDays  *[][]int32          //单个股票的周天
	MonthDays *[][]int32          //单个股票的月天
	YearDays  *[][]int32          //单个股票的年天
}

//所有股票
type SecurityList struct {
	Securitys *[]SingleSecurity
}

/*****************************************Functions***************************************/

//int型时间转Time类型（最小单位 天）
func IntToTime(date int) time.Time {
	swap := date % 10000
	year := date / 10000
	month := swap / 100
	day := swap % 100
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

//查询某一日期所在周的周日
func DateAdd(sid int32, date int) (time.Time, error) {
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
		basedate = fmt.Sprintf("%d%s", 24*6, "h")

	} else if strings.EqualFold(weekday, "Tuesday") {
		basedate = fmt.Sprintf("%d%s", 24*5, "h")

	} else if strings.EqualFold(weekday, "Wednesday") {
		basedate = fmt.Sprintf("%d%s", 24*4, "h")

	} else if strings.EqualFold(weekday, "Thursday") {
		basedate = fmt.Sprintf("%d%s", 24*3, "h")

	} else if strings.EqualFold(weekday, "Friday") {
		basedate = fmt.Sprintf("%d%s", 24*2, "h")

	} else if strings.EqualFold(weekday, "Saturday") {
		basedate = fmt.Sprintf("%d%s", 24*1, "h")

	} else {
		logging.Error("SID:%v------Invalid trade date...%v", sid, date)
		return sat, errors.New("周日有交易..")
	}

	dd, _ := time.ParseDuration(basedate)
	sat = baseTime.Add(dd) //Saturday（星期日）
	return sat, nil
}

//int型时间转Time类型（最小单位 月份）
func IntToMonth(date int) time.Time {
	swap := date % 10000
	year := date / 10000
	month := swap / 100

	return time.Date(year, time.Month(month), 0, 0, 0, 0, 0, time.UTC)
}

//K线数据写入相应文件的操作
func KlineWriteFile(sid int32, name string, data *[]byte) error {
	var filename string
	market := sid / 1000000
	if market == 100 {
		filename = fmt.Sprintf("%s/sh/%d/", cfg.File.Path, sid)
	} else if market == 200 {
		filename = fmt.Sprintf("%s/sz/%d/", cfg.File.Path, sid)
	} else {
		logging.Error("Monthline write file error...Invalid file path")
		return errors.New("Invalid file path")
	}

	err := os.MkdirAll(filename, 0777)
	if err != nil {
		fmt.Printf("%s", err)
	}

	err = ioutil.WriteFile(filename+name, *data, 0664)
	if err != nil {
		logging.Error("%v", err.Error())
		return err
	}
	return nil
}

func getDateToday() int32 {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	date, _ := strconv.Atoi(tm.Format("20060102"))
	return int32(date)
}
