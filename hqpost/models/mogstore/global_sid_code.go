//股票代码表
package mogstore

import (
	"gopkg.in/mgo.v2/bson"
	"haina.com/market/hqpost/models"
	"haina.com/share/logging"
	"haina.com/share/store/mongo"
)

type SecurityCode struct {
	SID int32 `bson:"nSID"`
}

func GetSecuritCodeyModel() *mongo.Model {
	return &mongo.Model{
		TableName: models.MOGON_SECURITY_TABLE,
	}
}

//股票代码表
func GetSecurityCodeTableFromMG() (*[]*SecurityCode, error) {
	var codes []*SecurityCode
	mogo := GetSecuritCodeyModel()

	exps := bson.M{}
	err := mogo.GetMulti(exps, &codes, 0, "nSID")
	if err != nil {
		logging.Error("%v", err.Error())
	}
	logging.Debug("lenght of sidcode tables:%v", len(codes))

	return &codes, err
}
