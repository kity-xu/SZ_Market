//市场代码表
package tb_security

import (
	"haina.com/market/hqinit/servers"
)

type TagSecurityInfo struct {
	NSID        int32  `bson:"nSID"`
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

//市场代码表
func GetSecurityInfoTableFromMG() *[]*TagSecurityInfo {
	var secus []*TagSecurityInfo

	TagI := new(servers.TagSecurityInfo).GetStockInfo("s2")

	for _, ite := range TagI {
		var tsi TagSecurityInfo
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
		secus = append(secus, &tsi)
	}
	return &secus
}
