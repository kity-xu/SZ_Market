package main

import (
	"bytes"
	"encoding/binary"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"haina.com/market/hqtools/dklinetools/financemysql"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type DKLine struct {
	NSID     int32  // 证券ID
	NTime    int32  // 时间 unix time
	NPreCPx  int32  // 昨收价 * 10000
	NOpenPx  int32  // 开盘价 * 10000
	NHighPx  int32  // 最高价 * 10000
	NLowPx   int32  // 最低价 * 10000
	NLastPx  int32  // 最新价 * 10000
	LlVolume int64  // 成交量
	LlValue  int64  // 成交额 * 10000
	NAvgPx   uint32 // 平均价 * 10000

}

func main() {
	logging.SetLogModel(true, false)
	logging.Info("begin..")
	// 查询finance数据库历史日K线
	conn, err := dbr.Open("mysql", "publisher:Haina$A7Kha@tcp(123.56.30.141:3306)/finchinafcdd?charset=utf8", nil)
	if err != nil {
		logging.Debug("mysql onn", err)
	}
	sess := conn.NewSession(nil)

	// 获取所有沪深股票代码
	CCodes, err := new(financemysql.ComCode).GetComCodeList(sess)
	var stocks []financemysql.Stock
	var err1 error
	for _, item := range CCodes {
		// 根据证券代码查询历史K线
		logging.Info("item.secode %v", item.SECODE.String)
		stocks, err1 = new(financemysql.Stock).GetSKTListFC(sess, item.SECODE.String)
		if err1 != nil {
			logging.Info("K线历史 %v", err1)
			return
		}
		if stocks != nil {
			logging.Info("GO")

		}

		if item.EXCHANGE.String == "001002" {
			var addstr = "E:/hqdata/sh/100" + item.SYMBOL.String
			// linux 下路径  var addstr = "/home/hqdata/sh/100" + item.SYMBOL.String
			WriteFileInfo(addstr, stocks, "100"+item.SYMBOL.String)
		}
		// 001003 深圳交易市场
		if item.EXCHANGE.String == "001003" {
			var addstr = "E:/hqdata/sz/200" + item.SYMBOL.String
			// linux 下路径  var addstr = "/home/hqdata/sz/200" + item.SYMBOL.String
			WriteFileInfo(addstr, stocks, "200"+item.SYMBOL.String)
		}
	}

}

// 写入文件
func WriteFileInfo(add string, sto []financemysql.Stock, snid string) {
	// 检查目录如果没有创建
	lib.CheckDir(add)
	var adds = add + "//dk.dat"
	// linux 下路径  var adds = add +"/dk.dat"
	file, err := os.OpenFile(adds, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	logging.Info("详细信息 %v", sto)
	i, err := strconv.Atoi(snid)
	if err != nil {
		logging.Info("类型转换出错 %v", snid)
	}
	for _, v := range sto {
		var sj DKLine
		sj.NSID = int32(i) // 证券ID
		sj.NTime = int32(v.TRADEDATE.Int64)
		sj.NPreCPx = int32(v.LCLOSE.Float64 * 10000)
		sj.NOpenPx = int32(v.TOPEN.Float64 * 10000)
		sj.NHighPx = int32(v.THIGH.Float64 * 10000)
		sj.NLowPx = int32(v.TLOW.Float64 * 10000)
		sj.NLastPx = int32(v.LCLOSE.Float64 * 10000)
		sj.LlVolume = v.VOL.Int64
		sj.LlValue = int64(v.AMOUNT.Float64 * 10000)
		sj.NAvgPx = uint32(v.AVGPRICE.Float64 * 10000)

		binary.Write(buf, binary.LittleEndian, sj)
	}

	file.Write(buf.Bytes())

}
