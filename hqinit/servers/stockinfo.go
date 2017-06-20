package servers

import (
	"strconv"

	"github.com/gocraft/dbr"

	fcm "haina.com/market/hqinit/models/fcmysql"
	"haina.com/share/logging"
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

var codes []*TagSecurityInfo

// 获取股票信息返回
func (this *TagSecurityInfo) GetStockInfo(sty string) []*TagSecurityInfo {

	// 获取沪深股票信息
	logging.Info("stockinfo begin==")
	conn, err := dbr.Open("mysql", "finchina:finchina@tcp(114.55.105.11:3306)/finchina?charset=utf8", nil)
	// 服务器用
	//conn, err := dbr.Open("mysql", "finchina:finchina@tcp(127.0.0.1:3306)/finchina?charset=utf8", nil)
	if err != nil {
		logging.Debug("mysql onn", err)
	}
	sess := conn.NewSession(nil)

	// s1 个股+指数   s2 个股  s3 指数
	if sty == "s1" {
		codes = nil
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

	}
	if sty == "s2" {
		codes = nil
		secNm3, err := new(fcm.FcSecuNameTab).GetSecuNmList(sess)
		if err != nil {
			logging.Info("个股查询finance出错 %v", err)
		}
		// 处理个股
		TreatingData(secNm3)
	}
	if sty == "s3" {
		codes = nil
		secNm4, err := new(fcm.FcSecuNameTab).GetExponentList(sess)
		if err != nil {
			logging.Info("指数查询finance出错 %v", err)
		}
		// 指数处理
		TreatingData(secNm4)
	}
	logging.Info("stockinfo End==")
	return codes
}

//  处理数据插入mongoDB
func TreatingData(secNm []*fcm.FcSecuNameTab) {

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

		codes = append(codes, &tsi)
	}

}
