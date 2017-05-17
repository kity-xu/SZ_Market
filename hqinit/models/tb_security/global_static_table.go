//证券静态数据
package tb_security

import (
	"gopkg.in/mgo.v2/bson"
	"haina.com/market/hqinit/models"
	"haina.com/share/logging"
	"haina.com/share/store/mongo"
)

type TagStockStatic struct {
	NSID           int32  `bson:"nSID"`           // 证券ID
	SzSType        string `bson:"szSType"`        // 证券类型
	SzStatus       string `bson:"szStatus"`       // 证券状态
	NListDate      int32  `bson:"nListDate"`      // 上市日期
	NLastTradeDate int32  `bson:"nLastTradeDate"` // 最近正常交易日期
	NDelistDate    int32  `bson:"nDelistDate"`    // 退市日期

	//NPreClose int32 `bson:"nPreClose"` // 前收价(*10000)
	//NHighLimitPx int32 `bson:"nHighLimitPx"` // 涨停价格(*10000)    // 作废
	//NLowLimitPx  int32 `bson:"nLowLimitPx"`  // 跌停价格(*10000)    // 作废

	LlCircuShare      int64   `bson:"llCircuShare"`      // 流通盘
	LlTotalShare      int64   `bson:"llTotalShare"`      // 总股本
	LlLast5Volume     int64   `bson:"llLast5Volume"`     // 最近5日成交总量(股)
	NEPS              float64 `bson:"nEPS"`              // 每股收益
	LlTotalProperty   float64 `bson:"llTotalProperty"`   // 总资产
	LlFlowProperty    int64   `bson:"llFlowProperty"`    // 流动资产
	NAVPS             int32   `bson:"nAVPS"`             // 每股净值
	LlMainIncoming    int64   `bson:"llMainIncoming"`    // 主营业务收入
	LlMainProfit      int64   `bson:"llMainProfit"`      // 主营业务利润
	LlTotalProfit     int64   `bson:"llTotalProfit"`     // 利润总额
	LlNetProfit       int64   `bson:"llNetProfit"`       // 净利润
	NHolders          int32   `bson:"nHolders"`          // 股东总数
	NReportDate       int32   `bson:"nReportDate"`       // 发布日期
	NCurrentRatio     int32   `bson:"nCurrentRatio"`     // 流通比率
	NQuickMovingRatio int32   `bson:"nQuickMovingRatio"` // 速动比例
}

func GetStockStaticModel() *mongo.Model {
	return &mongo.Model{
		TableName: models.MOGON_STATIC_TABLE,
	}
}

//股票代码表
func GetSecurityStaticTableFromMG() (*[]*TagStockStatic, error) {
	var tags []*TagStockStatic
	mogo := GetStockStaticModel()

	exps := bson.M{}
	err := mogo.GetMulti(exps, &tags, 0, "nSID")
	if err != nil {
		logging.Error("%v", err.Error())
	}
	logging.Debug("lenght of security static data tables:%v", len(tags))

	return &tags, err
}
