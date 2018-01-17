package main

import (
	"ProtocolBuffer/format/kline"
	"io/ioutil"

	"strconv"

	"github.com/golang/protobuf/proto"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"haina.com/market/hqtools/dklinetools/financemysql"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

var cnum chan int

func main() {
	logging.Info("begin..")
	// 查询finance数据库历史日K线
	//conn, err := dbr.Open("mysql", "finchina:finchina@tcp(114.55.105.11:3306)/finchina?charset=utf8", nil)
	// 服务器用
	//conn, err := dbr.Open("mysql", "finchina:finchina@tcp(172.16.1.60:3306)/finchina?charset=utf8", nil)
	//onn, err := dbr.Open("mysql", "f10_user:907fd7f123adcedd26a1bf537268f492@tcp(103.235.250.85:3306)/finchina?charset=utf8", nil)
	conn, err := dbr.Open("mysql", "f10_user:907fd7f123adcedd26a1bf537268f492@tcp(103.235.250.85:3306)/finchina?charset=utf8", nil)
	if err != nil {
		logging.Debug("mysql onn", err)
	}
	sess := conn.NewSession(nil)

	// 获取所有沪深股票代码
	CCodes, err := new(financemysql.ComCode).GetComCodeList(sess)

	if err != nil {
		logging.Error("select stock info error %v", err)
	}

	intx := len(CCodes) / 10
	for i := 0; i < 10; i++ {
		if i == 0 {
			go Compute(i, (i+1)*intx, CCodes, sess)
		} else if i == 9 {
			go Compute(i*intx+1, len(CCodes)-1, CCodes, sess)
		} else {
			go Compute(i*intx+1, (i+1)*intx, CCodes, sess)
		}
	}
	// 下面这个for循环的意义就是利用信道的阻塞，一直从信道里取数据，直到取得跟并发数一样的个数的数据，则视为所有goroutines完成。

	cnum = make(chan int, 10)
	for i := 0; i < 10; i++ {
		<-cnum
	}

	logging.Info("end.........The-successful-running.....")

}

func Compute(statidx, endidx int, ccom []financemysql.ComCode, sess *dbr.Session) {
	logging.Info("=========start:%v============end:%v", statidx, endidx)
	for i := statidx; i <= endidx; i++ {
		// 根据证券代码查询历史K线
		stocks, err1 := new(financemysql.Stock).GetSKTListFC(sess, ccom[i].SECODE.String)

		if err1 != nil {
			logging.Info("K线历史 %v", err1)
			return
		}

		if ccom[i].EXCHANGE.String == "001002" {
			var addstr = "/opt/develop/hgs/filestore/hqdata/sh/100" + ccom[i].SYMBOL.String
			//var addstr = "E:/hqdata/sh/100" + ccom[i].SYMBOL.String
			WriteFileInfo(addstr, stocks, "100"+ccom[i].SYMBOL.String)
		}
		// 001003 深圳交易市场
		if ccom[i].EXCHANGE.String == "001003" {
			var addstr = "/opt/develop/hgs/filestore/hqdata/sz/200" + ccom[i].SYMBOL.String
			//var addstr = "E:/hqdata/sz/200" + ccom[i].SYMBOL.String
			WriteFileInfo(addstr, stocks, "200"+ccom[i].SYMBOL.String)
		}
	}
	cnum <- 1 //goroutine结束时传送一个标示给信道。
}

// 写入文件
func WriteFileInfo(add string, sto []financemysql.Stock, snid string) {
	// 检查目录如果没有创建
	lib.CheckDir(add)
	if len(sto) < 1 {
		logging.Info("%v这支证券数据为空", snid)
		return
	}
	var adds = add + "/dk.dat"
	var klist kline.KInfoTable

	i, err := strconv.Atoi(snid)
	if err != nil {
		logging.Info("类型转换出错 %v", snid)
	}
	for _, v := range sto {
		if int(v.VOL.Int64) == 0 || v.AMOUNT.Float64 == 0.000 {
			//logging.Info("证券%v在交易日%v成交量为:%v成交额为:%v", i, v.TRADEDATE, v.VOL, v.AMOUNT)
			continue
		}
		var sj kline.KInfo
		sj.NSID = int32(i) // 证券ID
		sj.NTime = int32(v.TRADEDATE.Int64)
		sj.NPreCPx = int32(v.LCLOSE.Float64 * 10000)
		sj.NOpenPx = int32(v.TOPEN.Float64 * 10000)
		sj.NHighPx = int32(v.THIGH.Float64 * 10000)
		sj.NLowPx = int32(v.TLOW.Float64 * 10000)
		sj.NLastPx = int32(v.TCLOSE.Float64 * 10000)
		sj.LlVolume = v.VOL.Int64
		sj.LlValue = int64(v.AMOUNT.Float64 * 10000)
		sj.NAvgPx = uint32(v.AVGPRICE.Float64 * 10000)
		klist.List = append(klist.List, &sj)
	}
	data, err := proto.Marshal(&klist)
	if err != nil {
		logging.Error("%v", err.Error())
		return
	}

	if err = ioutil.WriteFile(adds, data, 0666); err != nil {
		logging.Error("%v", err.Error())
		return
	}
}
