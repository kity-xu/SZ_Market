//指数
package tb_security

import (
	"haina.com/market/hqinit/servers"
)

type StockIndexInfo struct {
	NSID        int32  `bson:"nSID"`        //股票指数SID会出现CN2008这种样式[8]byte
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

func GetStockIndexTableFromMG() *[]*StockIndexInfo {
	var table []*StockIndexInfo

	TagIk := new(servers.TagSecurityInfo).GetStockInfo("s3")
	for _, ite := range TagIk {
		var tsi StockIndexInfo
		tsi.NSID = ite.NSID
		tsi.NMarket = ite.NMarket
		tsi.SzSType = ite.SzSType
		tsi.SzStatus = ite.SzStatus
		tsi.SzSCode = ite.SzSCode
		tsi.SzSymbol = ite.SzSymbol
		tsi.SzISIN = ite.SzISIN
		tsi.SzSName = ite.SzSName
		tsi.SzSCName = ite.SzSCName
		tsi.SzDESC = ite.SzDESC
		tsi.SzPhonetic = ite.SzPhonetic
		tsi.SzCUR = ite.SzCUR
		tsi.SzIndusCode = ite.SzIndusCode
		table = append(table, &tsi)
	}
	return &table
}
