//市场代码表
package tb_security

import (
	"gopkg.in/mgo.v2/bson"
	"haina.com/market/hqpost/models"
	"haina.com/share/logging"
	"haina.com/share/store/mongo"
)

type TagSecurityInfo struct {
	//SID int32 `bson:"nSID"`
	NSID        int32  `bson:"nsid"`
	NMarket     int32  `bson:"nmarket"`     // 市场类型
	SzSType     string `bson:"szstype"`     // 证券类型
	SzStatus    string `bson:"szstatus"`    // 证券状态
	SzSCode     string `bson:"szscode"`     // 证券代码: 600036.SH
	SzSymbol    string `bson:"szsymbol"`    // 证券原始: 600036
	SzISIN      string `bson:"szisin"`      // 证券国际代码信息
	SzSName     string `bson:"szsname"`     // 证券名称 (超过24字节部分被省略)
	SzSCName    string `bson:"szscname"`    // 证券简体中文名称 (美股、港股超过40字节部分被省略
	SzDESC      string `bson:"szdesc"`      // 英文简称
	SzPhonetic  string `bson:"szphonetic"`  // 拼音
	SzCUR       string `bson:"szcur"`       // 币种
	SzIndusCode string `bson:"szinduscode"` // 行业代码
}

func GetMarketInfoModel() *mongo.Model {
	return &mongo.Model{
		TableName: models.MOGON_MARKET_TABLE,
	}
}

//市场代码表
func SZ_SecurityInfoTableFromMG() (*[]TagSecurityInfo, error) {
	var secus []TagSecurityInfo
	mogo := GetMarketInfoModel()

	exps := bson.M{
		"nmarket": 200000000, //深圳
	}
	err := mogo.GetMulti(exps, &secus, 0, "nsid")
	if err != nil {
		logging.Error("%v", err.Error())
	}
	logging.Debug("lenght of security market tables:%v", len(secus))

	return &secus, err
}

//市场代码表
func SH_SecurityInfoTableFromMG() (*[]TagSecurityInfo, error) {
	var secus []TagSecurityInfo
	mogo := GetMarketInfoModel()

	exps := bson.M{
		"nmarket": 100000000, //上海
	}
	err := mogo.GetMulti(exps, &secus, 0, "nsid")
	if err != nil {
		logging.Error("%v", err.Error())
	}
	logging.Debug("lenght of security market tables:%v", len(secus))

	return &secus, err
}
