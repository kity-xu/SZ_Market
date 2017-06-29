//证券静态数据
package tb_security

import (
	"haina.com/market/hqinit/servers"
	"haina.com/share/logging"
)

type TagStockStatic struct {
	NSID              int32  `bson:"nSID"`              // 证券ID
	SzSType           string `bson:"szSType"`           // 证券类型
	SzStatus          string `bson:"szStatus"`          // 证券状态
	NListDate         int32  `bson:"nListDate"`         // 上市日期
	NLastTradeDate    int32  `bson:"nLastTradeDate"`    // 最近正常交易日期
	NDelistDate       int32  `bson:"nDelistDate"`       // 退市日期
	LlCircuShare      int64  `bson:"llCircuShare"`      // 流通盘
	LlTotalShare      int64  `bson:"llTotalShare"`      // 总股本
	LlLast5Volume     int64  `bson:"llLast5Volume"`     // 最近5日成交总量(股)
	NEPS              int32  `bson:"nEPS"`              // 每股收益
	LlTotalProperty   int64  `bson:"llTotalProperty"`   // 总资产
	LlFlowProperty    int64  `bson:"llFlowProperty"`    // 流动资产
	NAVPS             int32  `bson:"nAVPS"`             // 每股净值
	LlMainIncoming    int64  `bson:"llMainIncoming"`    // 主营业务收入
	LlMainProfit      int64  `bson:"llMainProfit"`      // 主营业务利润
	LlTotalProfit     int64  `bson:"llTotalProfit"`     // 利润总额
	LlNetProfit       int64  `bson:"llNetProfit"`       // 净利润
	NHolders          int32  `bson:"nHolders"`          // 股东总数
	NReportDate       int32  `bson:"nReportDate"`       // 发布日期
	NCurrentRatio     int32  `bson:"nCurrentRatio"`     // 流通比例
	NQuickMovingRatio int32  `bson:"nQuickMovingRatio"` // 速动比例

	// 后加
	NEUndisProfit      int32 `bson:"nEUndisProfit"`      // 每股未分配利润
	NFlowLiab          int64 `bson:"nFlowLiab"`          // 流动负债
	NTotalLiabilities  int64 `bson:"nTotalLiabilities"`  // 负债总计
	NTotalHolderEquity int64 `bson:"nTotalHolderEquity"` // 股东权益合计
	NCapitalReserve    int64 `bson:"nCapitalReserve"`    // 资本公积金
	NIncomeInvestments int64 `bson:"nIncomeInvestments"` // 投资收益
}

//股票代码表
func GetSecurityStaticTableFromMG() *[]*TagStockStatic {
	var tags []*TagStockStatic

	stat := new(servers.TagStockStatic).GetStaticDataList()
	for _, ite := range stat {
		var tssc TagStockStatic

		tssc.NSID = ite.NSID
		tssc.SzSType = ite.SzSType
		tssc.SzStatus = ite.SzStatus
		tssc.NListDate = ite.NListDate
		tssc.NLastTradeDate = ite.NLastTradeDate
		tssc.NDelistDate = ite.NDelistDate
		tssc.LlCircuShare = ite.LlCircuShare
		tssc.LlTotalShare = ite.LlTotalShare
		tssc.LlLast5Volume = ite.LlLast5Volume
		tssc.NEPS = ite.NEPS
		tssc.LlTotalProperty = ite.LlTotalProperty
		tssc.LlFlowProperty = ite.LlFlowProperty
		tssc.NAVPS = ite.NAVPS
		tssc.LlMainIncoming = ite.LlMainIncoming
		tssc.LlMainProfit = ite.LlMainProfit
		tssc.LlTotalProfit = ite.LlTotalProfit
		tssc.LlNetProfit = ite.LlNetProfit
		tssc.NHolders = ite.NHolders
		tssc.NReportDate = ite.NReportDate
		tssc.NCurrentRatio = ite.NCurrentRatio
		tssc.NQuickMovingRatio = ite.NQuickMovingRatio

		tssc.NEUndisProfit = ite.NEUndisProfit
		tssc.NFlowLiab = ite.NFlowLiab
		tssc.NTotalLiabilities = ite.NTotalLiabilities
		tssc.NTotalHolderEquity = ite.NTotalHolderEquity
		tssc.NCapitalReserve = ite.NCapitalReserve
		tssc.NIncomeInvestments = ite.NIncomeInvestments

		tags = append(tags, &tssc)
	}
	logging.Debug("lenght of security static data tables:%v", len(tags))

	return &tags
}
