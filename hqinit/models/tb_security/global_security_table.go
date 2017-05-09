//市场代码表
package tb_security

import (
	"gopkg.in/mgo.v2/bson"
	"haina.com/market/hqinit/models"
	"haina.com/share/logging"
	"haina.com/share/store/mongo"
)

type TagSecurityInfo struct {
	//SID int32 `bson:"nSID"`
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

func GetMarketInfoModel() *mongo.Model {
	return &mongo.Model{
		TableName: models.MOGON_MARKET_TABLE,
	}
}

//市场代码表
func GetSecurityInfoTableFromMG() (*[]TagSecurityInfo, error) {
	var secus []TagSecurityInfo
	mogo := GetMarketInfoModel()

	exps := bson.M{
		"szSType": "101", //A股
	}
	err := mogo.GetMulti(exps, &secus, 0, "nSID")
	if err != nil {
		logging.Error("%v", err.Error())
	}
	logging.Debug("lenght of security market(include SH、SZ) tables:%v", len(secus))

	return &secus, err
}
