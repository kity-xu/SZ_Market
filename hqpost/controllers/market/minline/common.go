package minline

import (
	"ProtocolBuffer/format/kline"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"haina.com/share/lib"

	"io/ioutil"

	"haina.com/share/logging"

	"haina.com/market/hqpost/config"
)

var cfg *config.AppConfig

//历史分钟线
const (
	REDISKEY_SECURITY_MIN    = "hq:st:min:%d"    ///<证券分钟线数据(参数：sid) (calc写入)
	REDISKEY_SECURITY_HMIN   = "hq:st:hmin:%d"   ///<<证券历史分钟线数据(参数：sid) (hq-post写入)
	REDISKEY_SECURITY_HMIN5  = "hq:st:hmin5:%d"  ///<证券5分钟K线(参数：sid)
	REDISKEY_SECURITY_HMIN15 = "hq:st:hmin15:%d" ///<证券15分钟K线(参数：sid)
	REDISKEY_SECURITY_HMIN30 = "hq:st:hmin30:%d" ///<证券30分钟K线(参数：sid)
	REDISKEY_SECURITY_HMIN60 = "hq:st:hmin60:%d" ///<证券60分钟K线(参数：sid)
)

const (
	MIN_TOTAL = 241
	MIN_START = 930
	MIN_END   = 1500
)

type MinKline struct {
	sids *[]int32
	list AllMinLine
}

//个股
type SingleMin struct {
	Sid     int32                 //股票SID
	Time    []int32               //单个股票的历史日期
	Min     map[int32]kline.KInfo //单个股票的当天分钟数据
	Time_5  *[][]int32
	Time_15 *[][]int32
	Time_30 *[][]int32
	Time_60 *[][]int32
}

//所有股
type AllMinLine struct {
	All *[]*SingleMin
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

func KlineReadFile(sid int32, name string) ([]byte, error) {
	var filename string
	market := sid / 1000000
	if market == 100 {
		filename = fmt.Sprintf("%s/sh/%d/", cfg.File.Path, sid)
	} else if market == 200 {
		filename = fmt.Sprintf("%s/sz/%d/", cfg.File.Path, sid)
	} else {
		logging.Error("Monthline write file error...Invalid file path")
		return nil, errors.New("Invalid file path")
	}
	if !lib.IsFileExist(filename + name) {
		return nil, nil
	}
	return ioutil.ReadFile(filename + name)
}

func GetDateToday() int32 {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	date, _ := strconv.Atoi(tm.Format("20060102"))
	return int32(date)
}
