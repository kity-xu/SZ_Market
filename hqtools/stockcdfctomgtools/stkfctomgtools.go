package main

import (
	"time"

	"strconv"

	"github.com/gocraft/dbr"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	fcm "haina.com/market/hqtools/stockcdfctomgtools/financemysql"
	"haina.com/share/logging"
)

const (

	// 服务器mongodb://user:passwd@ip:port/db_name
	URL = "127.0.0.1:57017/hgs"
	//本地用
	//URL                  = "192.168.18.200:27017"
	GLOBAL_SECRITY_TABLE = "basic_securityinfo_table" // 证券静态数据monogoDb库
)

type TagSecurityInfo struct {
	NSID        int32  `bson:"nSID"`        // 证券ID
	NMarket     int32  `bson:"nMarket"`     // 市场类型
	SzSType     string `bson:"szSType"`     // 证券类型
	SzStatus    string `bson:"szStatus"`    // 证券状态
	SzSCode     string `bson:"szSCode"`     // 证券代码: 600036.SH
	SzSymbol    string `bson:"szSymbol"`    // 证券原始: 600036
	SzISIN      string `bson:"szISIN"`      // 证券国际代码信息
	SzSName     string `bson:"szSName"`     // 证券名称 (超过24字节部分被省略)
	SzSCName    string `bson:"szSCName"`    // 证券简体中文名称 (美股、港股超过40字节部分被省略
	SzDESC      string `bson:"szDESC"`      // 英文简称
	SzPhonetic  string `bson:"szPhonetic"`  // 拼音
	SzCUR       string `bson:"szCUR"`       // 币种
	SzIndusCode string `bson:"szIndusCode"` // 行业代码
}

// 读取沪深股票信息放到monogoDB中
func main() {

	// 获取沪深股票信息
	logging.Info("begin==")
	conn, err := dbr.Open("mysql", "finchina:finchina@tcp(114.55.105.11:3306)/finchina?charset=utf8", nil)
	// 服务器用
	// conn, err := dbr.Open("mysql", "finchina:finchina@tcp(172.16.1.60:3306)/finchina?charset=utf8", nil)
	if err != nil {
		logging.Debug("mysql onn", err)
	}
	sess := conn.NewSession(nil)
	secNm1, err := new(fcm.FcSecuNameTab).GetSecuNmList(sess)
	if err != nil {
		logging.Info("个股查询finance出错 %v", err)
	}
	// 处理个股
	TreatingData(secNm1)
	secNm2, err := new(fcm.FcSecuNameTab).GetExponentList(sess)
	if err != nil {
		logging.Info("指数查询finance出错 %v", err)
	}
	// 指数处理
	TreatingData(secNm2)
	logging.Info("End==")
}

//  处理数据插入mongoDB
func TreatingData(secNm []*fcm.FcSecuNameTab) {
	// monogoDB 插入

	mgo_conn, err := mgo.DialWithTimeout(URL, time.Second*10)
	if err != nil {
		logging.Info("monogoDB 插入出错 %v", err)
	}
	mgo_collection := mgo_conn.DB("hgs").C(GLOBAL_SECRITY_TABLE)
	//err = mgo_collection.Insert(&mgo_si) // 插入
	for _, item := range secNm {
		var tsi TagSecurityInfo
		swi := item.EXCHANGE.String
		sym := item.SYMBOL.String

		switch swi {
		case "001002":
			i, err := strconv.Atoi("100" + sym)
			tsi.NSID = int32(i)
			tsi.SzSCode = sym + ".SH"
			if err != nil {
				logging.Info("001002 sting 转 int 32 err %v", err)
			}
		case "001003":
			i, err := strconv.Atoi("200" + sym)
			tsi.NSID = int32(i)
			tsi.SzSCode = sym + ".SZ"
			if err != nil {
				logging.Info("001003 sting 转 int 32 err %v", err)
			}
		default:
			// 沪深以外的证券id
			i, err := strconv.Atoi("300" + sym)
			tsi.NSID = int32(i)
			tsi.SzSCode = sym + ".QT"
			if err != nil {
				logging.Info("qt sting 转 int 32 err %v", err)
			}
		}
		exh, err := strconv.Atoi(item.EXCHANGE.String)
		if err != nil {
			logging.Info("exchange sting 转 int 32 err %v", err)
		}
		switch exh {
		case 1002:
			tsi.NMarket = 100000000
		case 1003:
			tsi.NMarket = 200000000
		default:
			tsi.NMarket = 300000000
		}

		tsi.SzSType = item.SETYPE.String
		tsi.SzStatus = item.LISTSTATUS.String
		tsi.SzSymbol = item.SYMBOL.String
		tsi.SzISIN = item.SECURITYID.String
		tsi.SzSName = item.SENAME.String
		tsi.SzSCName = item.SESNAME.String
		tsi.SzDESC = item.SEENGNAME.String
		tsi.SzPhonetic = item.SESPELL.String
		tsi.SzCUR = item.CUR.String
		//tsi.SzIndusCode = item.
		//err = mgo_collection.Insert(&mgo_si) // 插入
		_, err = mgo_collection.Upsert(bson.M{"nSID": tsi.NSID}, &tsi) // nSID存在时更新，不存在时插入，此语句效率较慢
		if err != nil {
			logging.Info("insert error %v", err)
		}
	}
}
