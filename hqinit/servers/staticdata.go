package servers

import (
	"os"
	"strconv"
	"time"

	"github.com/LindsayBradford/go-dbf/godbf"
	_ "github.com/go-sql-driver/mysql"

	"haina.com/share/logging"

	"github.com/gocraft/dbr"
	stf "haina.com/market/hqinit/models/fcmysql"
	//	dkfm "haina.com/market/hqtools/dklinetools/financemysql"
	"bufio"
	"io"
	"strings"

	"haina.com/market/hqinit/config"

	"github.com/axgle/mahonia"
)

type SjsHqFile struct {
	NSID         int32   `bson:"nSID"`         // 证券ID		 <-- XXZQDM 证券代码
	SzStatus     string  `bson:"szStatus"`     // 证券状态    <-- XXJYZT 交易状态
	NListDate    int32   `bson:"nListDate"`    // 上市日期    <-- XXSSRQ 上市日期
	LlCircuShare int64   `bson:"llCircuShare"` // 流通盘      <-- XXLTGS 流通股数
	LlTotalShare int64   `bson:"llTotalShare"` // 总股本      <-- XXZFXL 总发行量
	NEPS         float64 `bson:"nEPS"`         // 每股收益    <-- XXSNLR 上年每股利润
	NAVPS        float64 `bson:"nAVPS"`        // 每股净值    <-- XXMGMZ 每股面值
}

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
	NCurrentRatio     int32  `bson:"nCurrentRatio"`     // 流动比率
	NQuickMovingRatio int32  `bson:"nQuickMovingRatio"` // 速动比率
	// 后加
	NEUndisProfit      int32 `bson:"nEUndisProfit"`      // 每股未分配利润
	NFlowLiab          int64 `bson:"nFlowLiab"`          // 流动负债
	NTotalLiabilities  int64 `bson:"nTotalLiabilities"`  // 负债总计
	NTotalHolderEquity int64 `bson:"nTotalHolderEquity"` // 股东权益合计
	NCapitalReserve    int64 `bson:"nCapitalReserve"`    // 资本公积金
	NIncomeInvestments int64 `bson:"nIncomeInvestments"` // 投资收益
	// 净利润处理给计算市盈率使用
	LlNetProfitS int64 `bson:"llNetProfitS"` // 净利润统计
}

// 整理静态数据放到monogoDB中
func (this *TagStockStatic) GetStaticDataList(cfg *config.AppConfig) []*TagStockStatic {

	tatic := StockTreatingData()

	// 把静态数据传到文件解析进行比对校验

	//return AnalysisFileUpMongodb(tatic, cfg)
	return tatic
}

// 处理个股静态数据
func StockTreatingData() []*TagStockStatic {
	var tsd []*TagStockStatic
	// 获取沪深股票信息
	secNm, err := stf.NewFcSecuNameTab().GetSecuNmList()

	if err != nil {
		logging.Info("查询finance出错 %v", err)
	}
	// 查询所有停牌股票
	ntrdule, err := stf.NewTQ_OA_NTRDSCHEDULE().GetNtrdsList()
	if err != nil {
		logging.Error("select TQ_OA_NTRDSCHEDULE error:%v", err)
	}
	// 遍历所有沪深股票
	for index, item := range secNm {
		var tss TagStockStatic
		// 根据证券id获取证券信息
		basinfo, err := stf.NewTQ_SK_BASICINFO().GetBasicinfoList(item.SYMBOL.String)

		if err != nil {
			if err == dbr.ErrNotFound {
				logging.Info("查询证券信息未找到数据 %v", err)
			} else {
				logging.Info("查询证券信息err %v", err)
			}
		}
		// 转换证券代码
		swi := item.EXCHANGE.String
		sym := item.SYMBOL.String
		switch swi {
		case "001002":
			i, err := strconv.Atoi("100" + sym)
			tss.NSID = int32(i)
			if err != nil {
				logging.Info("sting 转 int 32 err %v", err)
			}
		case "001003":
			i, err := strconv.Atoi("200" + sym)
			tss.NSID = int32(i)
			if err != nil {
				logging.Info("sting 转 int 32 err %v", err)
			}
		}
		tss.SzSType = item.SETYPE.String

		// 证券状态 处理停牌非停牌---------------------------------------------

		var isasis = false // 判断是佛有停牌
		for _, nitel := range ntrdule {
			// 有停牌
			if item.SECODE.String == nitel.SECODE.String {
				isasis = true
				break
			} else {
				isasis = false
			}
		}
		if isasis == true {
			tss.SzStatus = "0"
		}
		if isasis == false {
			tss.SzStatus = "1"
		}
		// ----------------------------------------------------------------

		lde, err := strconv.Atoi(basinfo.LISTDATE.String)
		tss.NListDate = int32(lde)
		dse, err := strconv.Atoi(basinfo.DELISTDATE.String)
		tss.NDelistDate = int32(dse)

		// 查询股本结构变化
		sharch, err := stf.NewTQ_SK_SHARESTRUCHG().GetSingleInfo(item.COMPCODE.String)

		if err != nil {
			logging.Info("查询股本结构 error")
		} else {

			tss.LlTotalShare = int64(sharch.TOTALSHARE.Float64 * 10000)
			tss.LlCircuShare = int64(sharch.CIRCSKAMT.Float64 * 10000)
		}

		// 根据公司内码获取股东信息
		shdn, errs := stf.NewTQ_SK_SHAREHOLDERNUM().GetSingleInfo(item.COMPCODE.String)

		if errs != nil {
			if err == dbr.ErrNotFound {
				logging.Info("查询股东信息未找到数据 %v", err)
			} else {
				logging.Info("查询股东信息err %v", errs)
			}
		} else {
			tost, err := strconv.Atoi(shdn.TOTALSHAMT.String)
			tss.NHolders = int32(tost)
			if err != nil {
				logging.Info("类型转换 error %v", err)
			}
		}

		// 根据公司内码查询股票历史信息   五日交易量
		ineqt, err := stf.NewTQ_SK_INTERVALQT().GetSingleInfo(item.SECODE.String)

		tss.NLastTradeDate = 0

		tss.LlLast5Volume = ineqt.VOL5D
		// 查询一般企业利润
		tspe, err := stf.NewTQ_FIN_PROINCSTATEMENTNEW().GetSingleInfo(item.COMPCODE.String)

		if err != nil {
			if err == dbr.ErrNotFound {
				logging.Info("查询一般企业利润未找到数据 %v", err)
			} else {
				logging.Info("查询一般企业利润出错 error %v", err)
			}
		}
		tss.NEPS = int32(tspe.BASICEPS.Float64 * 10000)
		tss.LlTotalProfit = int64(tspe.TOTPROFIT.Float64)
		tss.LlNetProfit = int64(tspe.PARENETP.Float64)
		// 净利润特殊处理------------------------------------------------------
		//		Proin, err := stf.NewTQ_FIN_PROINCSTATEMENTNEW().GetProinList(item.COMPCODE.String)
		//		if err != nil {
		//			logging.Error("seelct TQ_FIN_PROINCSTATEMENTNEW error")
		//		}
		//		// 获取当前年份
		//		year := time.Now().Year()
		//		var yint = 0

		//		for _, pite := range Proin {
		//			// 统计当前年份数据条数
		//			endint, err := strconv.Atoi(pite.ENDDATE.String[:4])
		//			if err != nil {
		//				logging.Error("ENDDATE type conversion error:%v", err)
		//			}
		//			profs := int64(pite.NETPROFIT.Float64)
		//			if endint == year {
		//				tss.LlNetProfitS += profs
		//				// 计算今年发布了几次
		//				yint += 1
		//			}
		//		}
		//		if len(Proin) == 5 {
		//			switch yint {
		//			case 0: // 如果今年没发布 取数组第 [4]
		//				tss.LlNetProfitS = int64(Proin[4].NETPROFIT.Float64)
		//			case 1: // 如果今年发布一个 取数组第 [4]+([3]-[0])
		//				val1 := int64(Proin[4].NETPROFIT.Float64 + (Proin[3].NETPROFIT.Float64 - Proin[0].NETPROFIT.Float64))
		//				tss.LlNetProfitS = val1
		//			case 2: // 如果今年发布二个 取数组第 [4]+([2]-[0])
		//				val2 := int64(Proin[4].NETPROFIT.Float64 + (Proin[2].NETPROFIT.Float64 - Proin[0].NETPROFIT.Float64))
		//				tss.LlNetProfitS = val2
		//			case 3: // 如果今年发布三个 取数组第 [4]+([1]-[0])
		//				val3 := int64(Proin[4].NETPROFIT.Float64 + (Proin[1].NETPROFIT.Float64 - Proin[0].NETPROFIT.Float64))
		//				tss.LlNetProfitS = val3
		//			case 4: //  今年数据 [4]
		//				tss.LlNetProfitS = int64(Proin[4].NETPROFIT.Float64)
		//			}
		//		}

		// ------------------------------------------------------------------
		endt, err := strconv.Atoi(tspe.ENDDATE.String)
		if err != nil {
			logging.Error("ENDDATE2 type conversion error:%v", err)
		}
		tss.NReportDate = int32(endt)

		tss.LlMainIncoming = int64(tspe.BIZINCO.Float64)
		tss.LlMainProfit = int64(tspe.PERPROFIT.Float64)
		tss.NIncomeInvestments = int64(tspe.INVEINCO.Float64)
		// 上市公司业绩快报 填充总资产  -- 改为取 负债表中数据
		//		tspce, err := stf.NewTQ_SK_PERFORMANCE().GetSingleInfo(item.COMPCODE.String)

		//		if err != nil {
		//			if err == dbr.ErrNotFound {
		//				logging.Info("上市公司业绩快报信息未找到数据 %v", err)
		//			} else {
		//				logging.Info("查询上市公司业绩快报信息出错 error %v", err)
		//			}
		//		} else {
		//			tss.NTotalHolderEquity = int64(tspce.TOTSHAREQUI.Float64 * 10000)
		//		}

		//查询一般企业资产负债信息 填充流动资产
		tfp, err := stf.NewTQ_FIN_PROBALSHEETNEW().GetSingleInfo(item.COMPCODE.String)

		if err != nil {
			if err == dbr.ErrNotFound {
				logging.Info("查询一般企业资产负债信息未找到数据 %v", err)
			} else {
				logging.Info("查询一般企业资产负债信息出错 error %v", err)
			}
		}
		tss.LlTotalProperty = int64(tfp.TOTASSET.Float64)
		tss.LlFlowProperty = int64(tfp.TOTCURRASSET.Float64)
		tss.NFlowLiab = int64(tfp.TOTALCURRLIAB.Float64)
		tss.NTotalLiabilities = int64(tfp.TOTLIAB.Float64)
		tss.NCapitalReserve = int64(tfp.CAPISURP.Float64)
		tss.NTotalHolderEquity = int64(tfp.RIGHAGGR.Float64 * 10000) // 股东权益
		// 每股净值= 归属母公司净利润/总股本
		if tss.LlTotalShare > 0 {
			tss.NAVPS = int32((tfp.PARESHARRIGH.Float64 / float64(tss.LlTotalShare)) * 10000)
		} else {
			tss.NAVPS = 0
		}

		//		profindex, err := stf.NewTQ_FIN_PROFINMAININDEX().GetSingleInfo(item.COMPCODE.String)
		//		if err != nil {
		//			logging.Error("select TQ_FIN_PROFINMAININDEX error :%v", err)
		//		}

		//		tss.NAVPS = int32(profindex.NAPS.Float64 * 10000)
		// 查询衍生财务指标信息 流动比率和速动比率
		tfpr, err := stf.NewTQ_FIN_PROINDICDATA().GetSingleInfo(item.COMPCODE.String)

		if err != nil {
			if err == dbr.ErrNotFound {
				logging.Info("查询衍生财务指标信未找到数据 %v", err)
			} else {
				logging.Info("查询衍生财务指标信息出错 error %v", err)
			}
		}
		tss.NCurrentRatio = int32(tfpr.CURRENTRT.Float64 * 10000)
		tss.NQuickMovingRatio = int32(tfpr.QUICKRT.Float64 * 10000)
		tss.NEUndisProfit = int32(tfpr.UPPS.Float64 * 10000)

		// 查询 一般企业利润 主营业务收入、利润、投资收益

		tsd = append(tsd, &tss)

		logging.Info("index ==%v=====%v", index, sym)
	}
	return tsd
}

// 分析 沪深市场证券基本信息文件修改 静态数据
func AnalysisFileUpMongodb(tss []*TagStockStatic, cfg *config.AppConfig) []*TagStockStatic {

	// 获取当前日期
	//timed := time.Now().Format("20060102")
	// 获取当前月日
	mod := time.Now().Format("0102")
	// 用来保存沪深所有个股
	sjshqfiles := []*SjsHqFile{}

	// 解析深证市场sjsxx.dbf文件 （证券信息库）
	//dbfTable, err := godbf.NewFromFile(cfg.File.Sjsxxdbfpath, "UTF8")
	// 服务器地址
	dbfTable, err := godbf.NewFromFile(cfg.File.Sjsxxdbfpath, "UTF8")
	if err != nil {
		logging.Info("==========%v", err)
		os.Exit(1)
	}
	for i := 0; i < dbfTable.NumberOfRecords(); i++ {
		var sjshqfile SjsHqFile
		// 获取第一列证券代码进行逻辑处理
		symstr, err := dbfTable.FieldValueByName(i, "XXZQDM")
		if err != nil {
			logging.Info("The XXZQDM column %v", err)
		}
		// 第一行
		if symstr == "000000" {
			continue
		}
		// 个股代码第一位为0或者3
		if symstr[:1] == "0" || symstr[:1] == "3" { // || symstr[:1] == "2" B基金
			// 个股代码第二位为0
			if symstr[1:2] == "0" {
				// 个股代码第三位为0或者2
				if symstr[2:3] >= "0" && symstr[2:3] <= "9" {
					nsid, err := strconv.Atoi("200" + symstr)
					if err != nil {
						logging.Info("这个%v证券代码转换in32 error", symstr)
					}
					sjshqfile.NSID = int32(nsid)
					//sjshqfile.SzStatus, err = dbfTable.FieldValueByName(i, "XXJYZT")
					listd, err := dbfTable.FieldValueByName(i, "XXSSRQ")
					if err != nil {
						logging.Info("这个%v证券代码解析上市日期 error", symstr)
					}
					listda, err := strconv.Atoi(listd)
					if err != nil {
						logging.Info("上市日期类型转换 error %v", err)
					}
					sjshqfile.NListDate = int32(listda)
					lcsstr, err := dbfTable.FieldValueByName(i, "XXLTGS")
					if err != nil {
						logging.Info("这个%v证券代码解析流通盘 error", symstr)
					}
					lcsint, err := strconv.Atoi(lcsstr)
					if err != nil {
						logging.Info("流通盘类型转换 error %v", err)
					}
					sjshqfile.LlCircuShare = int64(lcsint)
					ltsstr, err := dbfTable.FieldValueByName(i, "XXZFXL")
					if err != nil {
						logging.Info("这个%v证券代码解析总股本 error", symstr)
					}
					ltsint, err := strconv.Atoi(ltsstr)
					if err != nil {
						logging.Info("总股本类型转换 error %v", err)
					}
					sjshqfile.LlTotalShare = int64(ltsint)
					neps, err := dbfTable.FieldValueByName(i, "XXSNLR")
					if err != nil {
						logging.Info("这个%v证券代码解析每股收益 error", symstr)
					}
					nepsfl, err := strconv.ParseFloat(neps, 64)
					if err != nil {
						logging.Info("每股收益类型转换 error %v", err)
					}
					sjshqfile.NEPS = nepsfl
					navpsstr, err := dbfTable.FieldValueByName(i, "XXMGMZ")
					if err != nil {
						logging.Info("这个%v证券代码解析每股净值 error", symstr)
					}
					navpsint, err := strconv.ParseFloat(navpsstr, 64)
					if err != nil {
						logging.Info("每股净值类型转换 error %v", err)
					}
					sjshqfile.NAVPS = navpsint
					// 深交所数据
					sjshqfiles = append(sjshqfiles, &sjshqfile)
				}
			}
		}
	}

	// 上交所 证券处理
	//f, err := os.Open(cfg.File.Cpxxtxtpath + mod + ".txt") //打开文件
	// 服务器用
	f, err := os.Open(cfg.File.Cpxxtxtpath + mod + ".txt") //打开文件

	defer f.Close()                      //打开文件出错处理
	decoder := mahonia.NewDecoder("gbk") // 把原来ANSI格式的文本文件里的字符，用gbk进行解码。
	if nil == err {
		buff := bufio.NewReader(decoder.NewReader(f)) //读入缓存
		for {
			var sjshqfile SjsHqFile
			line, err := buff.ReadString('\n') //以'\n'为结束符读入一行
			if err != nil || io.EOF == err {
				logging.Info("reader ending!")
				break
			}
			//可以对一行进行处理

			strl := strings.Split(line, "|") // 根据|切割得到数组

			// ES 股票
			if strings.TrimSpace(strl[7]) == "ES" && strings.TrimSpace(strl[8]) == "ASH" {

				nsid, err := strconv.Atoi("100" + strings.TrimSpace(strl[0]))

				if err != nil {
					logging.Info("这个%v证券代码转换in32 error", strl[0])
				}
				// 解析上交所 cpxx0512文档 只有证券代码可以利用
				sjshqfile.NSID = int32(nsid)

			}
		}
	}
	// 沪深市场文档解析完成

	// 对比沪深数据
	for _, ite := range sjshqfiles {
		var isis = false
		for _, item := range tss {
			if ite.NSID == item.NSID {

				isis = true
				item.NSID = item.NSID
				item.SzSType = item.SzSType
				item.SzStatus = item.SzStatus

				item.NListDate = item.NListDate
				item.NLastTradeDate = item.NLastTradeDate
				item.NDelistDate = item.NDelistDate
				item.LlCircuShare = item.LlCircuShare
				item.LlTotalShare = item.LlTotalShare
				item.LlLast5Volume = item.LlLast5Volume
				//item.NEPS = int32(ite.NEPS * 10000)
				item.LlTotalProperty = item.LlTotalProperty
				item.LlFlowProperty = item.LlFlowProperty

				item.NAVPS = item.NAVPS * 10000
				item.NEPS = item.NEPS * 10000
				item.LlMainIncoming = item.LlMainIncoming
				item.LlMainProfit = item.LlMainProfit
				item.LlTotalProfit = item.LlTotalProfit
				item.LlNetProfit = item.LlNetProfit
				item.NHolders = item.NHolders
				item.NReportDate = item.NReportDate
				item.NCurrentRatio = item.NCurrentRatio
				item.NQuickMovingRatio = item.NQuickMovingRatio
				item.NEUndisProfit = item.NEUndisProfit
				item.NFlowLiab = item.NFlowLiab
				item.NTotalLiabilities = item.NTotalLiabilities
				item.NTotalHolderEquity = item.NTotalHolderEquity
				item.NCapitalReserve = item.NCapitalReserve
				item.NIncomeInvestments = item.NIncomeInvestments

			}
		}
		if isis == false {
			var tssc TagStockStatic
			tssc.NSID = ite.NSID
			//tssc.SzStatus = "1"
			tssc.SzSType = "101"
			tssc.NListDate = ite.NListDate
			tssc.LlCircuShare = ite.LlCircuShare
			//tssc.LlTotalShare = ite.LlTotalShare
			//tssc.NEPS = int32(ite.NEPS * 10000)
			//tssc.NAVPS = int32(ite.NAVPS * 10000)
			tss = append(tss, &tssc)
		}

	}
	return tss
}
