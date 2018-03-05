package security
//全市场品种静态文件 指数、股票
import (
	//"ProtocolBuffer/projects/hqinit/go/protocol"
	//"fmt"
	"strconv"

	//"github.com/golang/protobuf/proto"
	"haina.com/market/hqinit/config"


	//"haina.com/market/hqinit/models/tb_security"
	"haina.com/share/logging"
	//"haina.com/share/store/redis"
	"encoding/xml"
	"io/ioutil"
	"os"
	"haina.com/market/hqinit/servers"
)

type TagAllStockStatic struct {
	SzSCode     	  string  `xml:"Code"`// 证券代码: 600036.SH							len:SECURITY_CODE_LEN
	//NSID              int32   `xml:"ID"`// 证券ID
	SzSType           string  `xml:"Type"`// 证券类型
	SzStatus          string  `xml:"Status"`// 证券状态
	//NMarket     	  int32   `xml:"Market"`               // 市场类型
	NPreClosePx       int32   `xml:"PreClosePx"`               // 前收价(*10000)
	//SzSymbol    	  string  `xml:"Symbol"`// 证券原始: 600036								len:SECURITY_CODE_LEN
	SzISIN      	  string  `xml:"ISIN"`// 证券国际代码信息								len:SECURITY_ISIN_LEN
	SzSName     	  string  `xml:"Name"`// 证券名称 (超过24字节部分被省略)					len:SECURITY_NAME_LEN
	SzSCName    	  string  `xml:"CName"`// 证券简体中文名称 (美股、港股超过40字节部分被省略		len:SECURITY_NAME_LEN
	SzDESC      	  string  `xml:"DESC"`// 英文简称										len:SECURITY_DESC_LEN
	SzPhonetic  	  string  `xml:"Phonetic"`// 拼音											len:SECURITY_CODE_LEN
	SzCUR       	  string   `xml:"CUR"`              // 币种											len:4
	SzIndusCode 	  string  `xml:"IndusCode"`// 行业代码										len:INDUSTRY_CODE_LEN
	NListDate         int32   `xml:"ListDate"`// 上市日期
	NLastTradeDate    int32   `xml:"LastTradeDate"`// 最近正常交易日期
	NDelistDate       int32   `xml:"DelistDate"`// 退市日期
	LlCircuShare      int64   `xml:"CircuShare"`// 流通盘
	LlTotalShare      int64   `xml:"TotalShare"`// 总股本
	LlLast5Volume     int64   `xml:"Last5Volume"`// 最近5日成交总量(股)
	NEPS              int32   `xml:"EPS"`// 每股收益
	LlTotalProperty   int64   `xml:"TotalProperty"`// 总资产
	LlFlowProperty    int64   `xml:"FlowProperty"`// 流动资产
	NAVPS             int32   `xml:"AVPS"`// 每股净值
	LlMainIncoming    int64   `xml:"MainIncoming"`// 主营业务收入
	LlMainProfit      int64   `xml:"MainProfit"`// 主营业务利润
	LlTotalProfit     int64   `xml:"TotalProfit"`// 利润总额
	LlNetProfit       int64   `xml:"NetProfit"`// 净利润
	NHolders          int32   `xml:"Holders"`// 股东总数
	NReportDate       int32   `xml:"ReportDate"`// 发布日期
	NCurrentRatio     int32   `xml:"CurrentRatio"`// 流通比率
	NQuickMovingRatio int32   `xml:"QuickMovingRatio"`// 速动比例
	// 后加
	NEUndisProfit      int32  `xml:"EUndisProfit"`// 每股未分配利润
	NFlowLiab          int64  `xml:"CurrentLiability"`// 流动负债
	NTotalLiabilities  int64  `xml:"TotalLiabilities"`// 负债总计
	NTotalHolderEquity int64  `xml:"TotalHolderEquity"`// 股东权益合计
	NCapitalReserve    int64  `xml:"CapitalReserve"`// 资本公积金
	NIncomeInvestments int64  `xml:"IncomeInvestments"`// 投资收益
}
type Stype struct{
	Note string   `xml:"Note,attr"` // 说明
	T0 		string 		`xml:"T0"`
	T1 		string 		`xml:"T1"`
	T2 		string 		`xml:"T2"`
	T3 		string 		`xml:"T3"`
}
type SStatus struct{
	Note string   `xml:"Note,attr"` // 说明
	S0 		string 		`xml:"S0"`
	S1 		string 		`xml:"S1"`
	S2 		string 		`xml:"S2"`
	S3 		string 		`xml:"S3"`
}
type SecurityNote struct{
	Code 							string   `xml:"Code"`
	//ID 								string   `xml:"ID"`
	Type 							Stype   `xml:"Type"`
	Status 							SStatus   `xml:"Status"`
	//Market 							string   `xml:"Market"`
	PreClosePx 						string   `xml:"PreClosePx"`
	//Symbol 							string   `xml:"Symbol"`
	ISIN 							string   `xml:"ISIN"`
	Name 							string   `xml:"Name"`
	CName 							string   `xml:"CName"`
	DESC 							string   `xml:"DESC"`
	Phonetic 						string   `xml:"Phonetic"`
	CUR 							string   `xml:"CUR"`
	IndusCode 						string   `xml:"IndusCode"`
	ListDate 						string   `xml:"ListDate"`
	LastTradeDate 					string   `xml:"LastTradeDate"`
	DelistDate 						string   `xml:"DelistDate"`
	CircuShare 						string   `xml:"CircuShare"`
	TotalShare 						string   `xml:"TotalShare"`
	Last5Volume 					string   `xml:"Last5Volume"`
	EPS 							string   `xml:"EPS"`
	TotalProperty 					string   `xml:"TotalProperty"`
	FlowProperty 					string   `xml:"FlowProperty"`
	AVPS 							string   `xml:"AVPS"`
	MainIncoming 					string   `xml:"MainIncoming"`
	MainProfit 						string   `xml:"MainProfit"`
	TotalProfit 					string   `xml:"TotalProfit"`
	NetProfit 						string   `xml:"NetProfit"`
	Holders 						string   `xml:"Holders"`
	ReportDate 						string   `xml:"ReportDate"`
	CurrentRatio 					string   `xml:"CurrentRatio"`
	QuickMovingRatio 				string   `xml:"QuickMovingRatio"`
	EUndisProfit 					string   `xml:"EUndisProfit"`
	CurrentLiability 				string   `xml:"CurrentLiability"`
	TotalLiabilities 				string   `xml:"TotalLiabilities"`
	TotalHolderEquity 				string   `xml:"TotalHolderEquity"`
	CapitalReserve 					string   `xml:"CapitalReserve"`
	IncomeInvestments 				string   `xml:"IncomeInvestments"`
}

type Securities struct {
	StaticArr []TagAllStockStatic `xml:"Security"`
}
type root struct{
	SecurityNote  SecurityNote `xml:"SecurityNote"`
	Securities	  Securities	`xml:"Securities"`
}

//证券基本信息和单市场的证券代码表的实现
func UpdateAllSecurityStatic(cfg *config.AppConfig) {
	var (
		stype, status string
	)
	var root root
	//注释说明
	var senote SecurityNote
	var sstype Stype
	var sstatus SStatus
	sstype.Note = "证券类型"
	sstype.T0 = "0:未定义,1:沪市,2:深市"
	sstype.T1 = "I:指数,S:股票,F:基金,B:债券,W:权证,O:期权,U:期货,A:申购、收购、配号等证券业务,P:优先股,R:质押回购"
	sstype.T2 = "A:A股,B:B股...."
	sstype.T3 = "Z:A股主板,X:A股中小板,C:A股创业板..."

	sstatus.Note = "证券状态"
	sstatus.S0 = "-:未定义,0:正常,N:上市首日,R:恢复上市首日,D:退市,S:停牌,L:长期停牌,T:临时停牌"
	sstatus.S1 = "-:未定义,R:除权,D:除息,C:除权除息"
	sstatus.S2 = "-:未定义,*:*ST,S:ST,P:退市整理期,T:暂停上市后协议转让"
	sstatus.S3 = "-:未定义,L:债券投资者适当性要求类,G:未完成股改,R:公司再融资,S:增发股份上市,C:合约调整,V:网络投票,D:上网定价发行,J:上网竞价发行,F:国债挂牌分销"

	//senote.ID 				= "证券ID"
	senote.Type 			= sstype
	senote.Status 			= sstatus
	//senote.Market 			= "市场类型"
	senote.PreClosePx 		= "昨收价"
	senote.Code 			= "证券代码: 600036.SH"
	//senote.Symbol 			= "证券原始: 600036"
	senote.ISIN 			= "证券国际代码信息"
	senote.Name 			= "证券名称"
	senote.CName 			= "证券简体中文名称"
	senote.DESC 			= "英文简称"
	senote.Phonetic 		= "拼音"
	senote.CUR 				= "币种"
	senote.IndusCode 		= "行业代码"
	senote.ListDate 			= "上市日期"
	senote.LastTradeDate 			= "最近正常交易日期"
	senote.DelistDate 			= "退市日期"
	senote.CircuShare 			= "流通盘(股)"
	senote.TotalShare 			= "总股本(股)"
	senote.Last5Volume 			= "最近5日成交总量(股)"
	senote.EPS 					= "每股收益(*10000)"
	senote.TotalProperty 		= "总资产"
	senote.FlowProperty 		= "流动资产"
	senote.AVPS 				= "每股净值(*10000)"
	senote.MainIncoming 		= "主营业务收入"
	senote.MainProfit 			= "主营业务利润"
	senote.TotalProfit 			= "利润总额"
	senote.NetProfit 			= "净利润"
	senote.Holders 			= "股东总数"
	senote.ReportDate 			= "发布日期"
	senote.CurrentRatio 			= "流通比率(*10000)"
	senote.QuickMovingRatio 			= "速动比例(*10000)"
	senote.EUndisProfit 			= "每股未分配利润(*10000)"
	senote.CurrentLiability 			= "流动负债"
	senote.TotalLiabilities 			= "负债总计"
	senote.TotalHolderEquity 			= "股东权益合计(*10000)"
	senote.CapitalReserve 			= "资本公积金"
	senote.IncomeInvestments 			= "投资收益"
	//...TODO

	//securityname
	table := MarketTableExt()
	//static
	stable := getSecurityStaticExt(cfg)
	var staticmap map[int32]servers.TagStockStatic
	//转map
	staticmap = make(map[int32]servers.TagStockStatic)
	for _,v :=range stable{
		staticmap[v.NSID] = *v
	}

	var err error
	var staticxml []TagAllStockStatic
	var static Securities

	for _, v := range *table {
		stype, err = HainaSecurityType(strconv.Itoa(int(v.NSID)), v.SzSType)
		if err != nil {
			logging.Error("%v", err.Error())
		}
		status, err = HainaSecurityStatus(v.SzStatus)

		if err != nil {
			logging.Error("%v", err.Error())
		}

		val,ok:=staticmap[v.NSID]
		if(!ok && stype[1]=='S'){
			logging.Error("static info not found! nsid:%v", v.NSID)
		}

		biny := TagAllStockStatic{ //入文件的结构
			//NSID:              v.NSID,
			SzSType:           stype,
			SzStatus:          status,
			//NMarket:  		   v.NMarket,
			NPreClosePx:	   val.NPreClose,
			SzSCode: 	       v.SzSCode,
			//SzSymbol: 	       v.SzSymbol,
			SzISIN: 	       v.SzISIN,
			SzSName: 	       v.SzSName,
			SzSCName: 	   	   v.SzSCName,
			SzDESC: 	   	   v.SzDESC,
			SzPhonetic: 	   v.SzPhonetic,
			SzCUR: 	   		   v.SzCUR,
			SzIndusCode: 	   v.SzIndusCode,
			NListDate:         val.NListDate,
			NLastTradeDate:    val.NLastTradeDate,
			NDelistDate:       val.NDelistDate,
			LlCircuShare:      val.LlCircuShare,
			LlTotalShare:      val.LlTotalShare,
			LlLast5Volume:     val.LlLast5Volume,
			NEPS:              val.NEPS,
			LlTotalProperty:   val.LlTotalProperty,
			LlFlowProperty:    val.LlFlowProperty,
			NAVPS:             val.NAVPS,
			LlMainIncoming:    val.LlMainIncoming,
			LlMainProfit:      val.LlMainProfit,
			LlTotalProfit:     val.LlTotalProfit,
			LlNetProfit:       val.LlNetProfit,
			NHolders:          val.NHolders,
			NReportDate:       val.NReportDate,
			NCurrentRatio:     val.NCurrentRatio,
			NQuickMovingRatio: val.NQuickMovingRatio,

			NEUndisProfit:      val.NEUndisProfit,
			NFlowLiab:          val.NFlowLiab,
			NTotalLiabilities:  val.NTotalLiabilities,
			NTotalHolderEquity: val.NTotalHolderEquity,
			NCapitalReserve:    val.NCapitalReserve,
			NIncomeInvestments: val.NIncomeInvestments,
		}

		//
		staticxml = append(staticxml, biny)

		/*************************OVER******************************/


	}
	static.StaticArr = staticxml

	root.SecurityNote = senote
	root.Securities = static
	/*************************START******************************/
	//写xml
	output, err := xml.MarshalIndent(root, "  ", "    ")
	//加入XML头
	headerBytes := []byte(xml.Header)
	//拼接XML头和实际XML内容
	xmlOutPutData := append(headerBytes, output...)

	ioutil.WriteFile(cfg.File.SecurityStatic, xmlOutPutData, os.ModeAppend) // 服务器用

}

