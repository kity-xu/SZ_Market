package main

import (
	"bytes"
	"encoding/binary"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	dkl "haina.com/market/hqtools/dklinetools/financemysql"
	ilt "haina.com/market/hqtools/indexlinetools/financemysql"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type IndexDKLine struct {
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
	// 查询finance数据库历史指数K线数据
	conn, err := dbr.Open("mysql", "finchina:finchina@tcp(114.55.105.11:3306)/finchina?charset=utf8", nil)
	if err != nil {
		logging.Debug("mysql onn", err)
	}
	sess := conn.NewSession(nil)
	// 获取指数信息
	CCodes, err := new(dkl.ComCode).GetIndexInfoList(sess)

	var index []ilt.TQ_QT_INDEX
	var err1 error
	for _, item := range CCodes {
		// 根据指数历史K线
		//logging.Info("item.secode %v", item.SECODE.String)
		index, err1 = new(ilt.TQ_QT_INDEX).GetIndexInfoList(sess, item.SECODE.String)
		if err1 != nil {
			logging.Info("K线历史 %v", err1)
			return
		}

		if item.EXCHANGE.String == "001002" {
			var addstr = "E:/hqdata/sh/100" + item.SYMBOL.String
			// linux 下路径  var addstr = "/home/hqdata/sh/100" + item.SYMBOL.String
			WriteFileInfo(addstr, index, "100"+item.SYMBOL.String)
		}
		// 001003 深圳交易市场
		if item.EXCHANGE.String == "001003" {
			var addstr = "E:/hqdata/sz/200" + item.SYMBOL.String
			// linux 下路径  var addstr = "/home/hqdata/sz/200" + item.SYMBOL.String
			WriteFileInfo(addstr, index, "200"+item.SYMBOL.String)
		}
	}
	logging.Info("end..")
}

// 写入文件
func WriteFileInfo(add string, sto []ilt.TQ_QT_INDEX, snid string) {
	// 检查目录如果没有创建
	lib.CheckDir(add)
	if len(sto) < 1 {
		logging.Info("%v这支证券数据为空", snid)
		return
	}
	var adds = add + "//indexdk.dat"
	// linux 下路径  var adds = add +"/indexdk.dat"
	file, err := os.OpenFile(adds, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)

	i, err := strconv.Atoi(snid)
	if err != nil {
		logging.Info("类型转换出错 %v", snid)
	}
	for _, v := range sto {
		var sj IndexDKLine
		sj.NSID = int32(i) // 证券ID
		sj.NTime = int32(v.TRADEDATE.Float64)
		sj.NPreCPx = int32(v.LCLOSE.Float64 * 10000)
		sj.NOpenPx = int32(v.TOPEN.Float64 * 10000)
		sj.NHighPx = int32(v.THIGH.Float64 * 10000)
		sj.NLowPx = int32(v.TLOW.Float64 * 10000)
		sj.NLastPx = int32(v.LCLOSE.Float64 * 10000)
		sj.LlVolume = int64(v.VOL.Float64)
		sj.LlValue = int64(v.AMOUNT.Float64 * 10000)
		//sj.NAvgPx = uint32(v.AVGPRICE.Float64 * 10000)  指数表中没有平均价格

		binary.Write(buf, binary.LittleEndian, sj)
	}

	_, errs := file.Write(buf.Bytes())
	if errs != nil {
		logging.Info("write file %v", errs)
	}

}
