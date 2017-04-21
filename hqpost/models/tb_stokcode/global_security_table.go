package tb_stokcode

import (
	"gopkg.in/mgo.v2/bson"
	"haina.com/market/hqpost/models"
	"haina.com/share/logging"
	"haina.com/share/store/mongo"
)

type Code struct {
	//Id_    bson.ObjectId `bson:"_id"`
	SID int32 `bson:"nSID"`
	//SCName string        `bson:"szSCName"`
}

func GetMongoModel() *mongo.Model {
	return &mongo.Model{
		TableName: models.MOGON_SECURITY_TABLE,
	}
}

func GetSecurityTableFromMG() ([]Code, error) {
	var codes []Code
	mogo := GetMongoModel()

	exps := bson.M{}
	err := mogo.GetMulti(exps, &codes, 5000, "nSID")
	if err != nil {
		logging.Error("%v", err.Error())
	}
	logging.Debug("lenght of stockcode tables:%v", len(codes))

	return codes, err
}
