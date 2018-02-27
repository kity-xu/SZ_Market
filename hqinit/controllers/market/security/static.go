//证券静态数据
package security

import (
	//"ProtocolBuffer/projects/hqinit/go/protocol"
	"bytes"
	"encoding/binary"
	//"fmt"
	"strconv"

	//"github.com/golang/protobuf/proto"
	"haina.com/market/hqinit/config"
	"haina.com/share/lib"

	. "haina.com/market/hqinit/controllers"

	"haina.com/market/hqinit/models/tb_security"
	"haina.com/share/logging"
	//"haina.com/share/store/redis"
)

type TagStockStatic struct {
	NSID              int32   // 证券ID
	SzSType           [4]byte // 证券类型
	SzStatus          [4]byte // 证券状态
	NListDate         int32   // 上市日期
	NLastTradeDate    int32   // 最近正常交易日期
	NDelistDate       int32   // 退市日期
	LlCircuShare      int64   // 流通盘
	LlTotalShare      int64   // 总股本
	LlLast5Volume     int64   // 最近5日成交总量(股)
	NEPS              int32   // 每股收益
	LlTotalProperty   int64   // 总资产
	LlFlowProperty    int64   // 流动资产
	NAVPS             int32   // 每股净值
	LlMainIncoming    int64   // 主营业务收入
	LlMainProfit      int64   // 主营业务利润
	LlTotalProfit     int64   // 利润总额
	LlNetProfit       int64   // 净利润
	NHolders          int32   // 股东总数
	NReportDate       int32   // 发布日期
	NCurrentRatio     int32   // 流通比率
	NQuickMovingRatio int32   // 速动比例
	// 后加
	NEUndisProfit      int32 // 每股未分配利润
	NFlowLiab          int64 // 流动负债
	NTotalLiabilities  int64 // 负债总计
	NTotalHolderEquity int64 // 股东权益合计
	NCapitalReserve    int64 // 资本公积金
	NIncomeInvestments int64 // 投资收益
}

// 静态数据
func getSecurityStatic(cfg *config.AppConfig) *[]*tb_security.TagStockStatic {
	return tb_security.GetSecurityStaticTableFromMG(cfg)
}

func UpdateSecurityStaticInfo(cfg *config.AppConfig) {

	var (
		stype, status string
	)
	stable := getSecurityStatic(cfg)
	var err error
	//入文件
	buffer := new(bytes.Buffer)

	for _, v := range *stable {
		stype, err = HainaSecurityType(strconv.Itoa(int(v.NSID)), v.SzSType)
		if err != nil {
			logging.Error("%v", err.Error())
		}
		status, err = HainaSecurityStatus(v.SzStatus)
		if err != nil {
			logging.Error("%v", err.Error())
		}
		//tag := protocol.StockStatic{ //入redis的结构
		//	NSID:              v.NSID,
		//	SzSType:           stype,
		//	SzStatus:          status,
		//	NListDate:         v.NListDate,
		//	NLastTradeDate:    v.NLastTradeDate,
		//	NDelistDate:       v.NDelistDate,
		//	LlCircuShare:      v.LlCircuShare,
		//	LlTotalShare:      v.LlTotalShare,
		//	LlLast5Volume:     v.LlLast5Volume,
		//	NEPS:              v.NEPS,
		//	LlTotalProperty:   v.LlTotalProperty,
		//	LlFlowProperty:    v.LlFlowProperty,
		//	NAVPS:             v.NAVPS,
		//	LlMainIncoming:    v.LlMainIncoming,
		//	LlMainProfit:      v.LlMainProfit,
		//	LlTotalProfit:     v.LlTotalProfit,
		//	LlNetProfit:       v.LlNetProfit,
		//	NHolders:          v.NHolders,
		//	NReportDate:       v.NReportDate,
		//	NCurrentRatio:     v.NCurrentRatio,
		//	NQuickMovingRatio: v.NQuickMovingRatio,
		//
		//	NEUndisProfit:      v.NEUndisProfit,
		//	NFlowLiab:          v.NFlowLiab,
		//	NTotalLiabilities:  v.NTotalLiabilities,
		//	NTotalHolderEquity: v.NTotalHolderEquity,
		//	NCapitalReserve:    v.NCapitalReserve,
		//	NIncomeInvestments: v.NIncomeInvestments,
		//}
		biny := TagStockStatic{ //入文件的结构
			NSID:              v.NSID,
			SzSType:           StringToByte_4(stype),
			SzStatus:          StringToByte_4(status),
			NListDate:         v.NListDate,
			NLastTradeDate:    v.NLastTradeDate,
			NDelistDate:       v.NDelistDate,
			LlCircuShare:      v.LlCircuShare,
			LlTotalShare:      v.LlTotalShare,
			LlLast5Volume:     v.LlLast5Volume,
			NEPS:              v.NEPS,
			LlTotalProperty:   v.LlTotalProperty,
			LlFlowProperty:    v.LlFlowProperty,
			NAVPS:             v.NAVPS,
			LlMainIncoming:    v.LlMainIncoming,
			LlMainProfit:      v.LlMainProfit,
			LlTotalProfit:     v.LlTotalProfit,
			LlNetProfit:       v.LlNetProfit,
			NHolders:          v.NHolders,
			NReportDate:       v.NReportDate,
			NCurrentRatio:     v.NCurrentRatio,
			NQuickMovingRatio: v.NQuickMovingRatio,

			NEUndisProfit:      v.NEUndisProfit,
			NFlowLiab:          v.NFlowLiab,
			NTotalLiabilities:  v.NTotalLiabilities,
			NTotalHolderEquity: v.NTotalHolderEquity,
			NCapitalReserve:    v.NCapitalReserve,
			NIncomeInvestments: v.NIncomeInvestments,
		}

		//转PB
		//data, err := proto.Marshal(&tag)
		//if err != nil {
		//	logging.Error("Encode protocbuf of week Line error...%v", err.Error())
		//	return
		//}

		//入redis
		//key := fmt.Sprintf(REDISKEY_SECURITY_STATIC, tag.NSID)
		//if err := redis.Set(key, data); err != nil {
		//	logging.Fatal("%v", err)
		//}

		//缓冲二进制数据
		if err := binary.Write(buffer, binary.LittleEndian, &biny); err != nil {
			logging.Fatal(err)
		}
	}
	//入文件
	lib.CheckDir(cfg.File.Path)
	file, err := OpenFile(cfg.File.Path + cfg.File.StaticName)
	if err != nil {
		return
	}


	_, err1 := file.Write(buffer.Bytes())
	if err1 != nil {
		logging.Error("Write file error...")
	}

	defer file.Close()
}
