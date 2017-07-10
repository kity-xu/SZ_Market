package servers

import (
	"strconv"

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

// 查询当日新股
type BasicInfoN struct {
	LISTDATE string // 上市日期
	SYMBOL   string // 证券代码
}

// 获取当日新股
func (this *BasicInfoN) GetBasiN() []*BasicInfoN {
	nbi, err := fcm.NewTQ_SK_BASICINFO().GetNewBasicinfo()
	if err != nil {
		logging.Info("查询当日新股 error %v", err)
	}
	var binl []*BasicInfoN
	for _, itb := range nbi {
		var bi BasicInfoN
		bi.LISTDATE = itb.LISTDATE.String
		bi.SYMBOL = itb.SYMBOL.String
		binl = append(binl, &bi)
	}
	return binl
}

var codes []*TagSecurityInfo

// 获取股票信息返回
func (this *TagSecurityInfo) GetStockInfo(sty string) []*TagSecurityInfo {

	// s1 个股+指数   s2 个股  s3 指数
	if sty == "s1" {
		codes = nil
		secNm1, err := fcm.NewFcSecuNameTab().GetSecuNmList()
		if err != nil {
			logging.Info("个股查询finance出错 %v", err)
		}
		// 处理个股
		TreatingData(secNm1)
		secNm2, err := fcm.NewFcSecuNameTab().GetExponentList()
		if err != nil {
			logging.Info("指数查询finance出错 %v", err)
		}
		// 指数处理
		TreatingData(secNm2)

	}
	if sty == "s2" {
		codes = nil
		secNm3, err := fcm.NewFcSecuNameTab().GetSecuNmList()
		if err != nil {
			logging.Info("个股查询finance出错 %v", err)
		}
		// 处理个股
		TreatingData(secNm3)
	}
	if sty == "s3" {
		codes = nil
		secNm4, err := fcm.NewFcSecuNameTab().GetExponentList()
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

	bas, err := fcm.NewTQ_SK_BASICINFO().GetNewBasicinfo()
	if err != nil {
		logging.Info("查询当日新股error %v", err)
	}
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

		// 如果当日有新股 新股名字加N
		var syisis = false
		if len(bas) > 0 {
			for _, ibas := range bas {
				if item.SYMBOL.String == ibas.SYMBOL.String {
					syisis = true
				}
			}
		}
		if syisis == true {
			rsn := []rune(item.SENAME.String)
			rszs := []rune(item.SESNAME.String)
			tsi.SzSName = "N" + string(rsn[0]) + string(rszs[1])
			tsi.SzSCName = "N" + string(rszs[0]) + string(rszs[1])
			// 新股拼音前加N
			tsi.SzPhonetic = "N" + item.SESPELL.String[:1] + item.SESPELL.String[1:2]
		} else {
			tsi.SzPhonetic = item.SESPELL.String
			tsi.SzSName = item.SENAME.String
			tsi.SzSCName = item.SESNAME.String
		}

		tsi.SzDESC = item.SEENGNAME.String

		tsi.SzCUR = item.CUR.String
		//tsi.SzIndusCode = item.

		codes = append(codes, &tsi)
	}

}
