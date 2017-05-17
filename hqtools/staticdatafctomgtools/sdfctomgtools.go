package main

import (
	"strconv"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gocraft/dbr"
	dkfm "haina.com/market/hqtools/dklinetools/financemysql"
	stf "haina.com/market/hqtools/staticdatafctomgtools/financemysql"
	fms "haina.com/market/hqtools/stockcdfctomgtools/financemysql"
	"haina.com/share/logging"
)

const (
	URL                  = "192.168.18.200:27017"
	GLOBAL_SECRITY_TABLE = "basic_staticdata_table_test" // 证券静态数据monogoDb库
)

type TagStockStatic struct {
	NSID              int32   `bson:"nSID"`              // 证券ID
	SzSType           string  `bson:"szSType"`           // 证券类型
	SzStatus          string  `bson:"szStatus"`          // 证券状态
	NListDate         int32   `bson:"nListDate"`         // 上市日期
	NLastTradeDate    int32   `bson:"nLastTradeDate"`    // 最近正常交易日期
	NDelistDate       int32   `bson:"nDelistDate"`       // 退市日期
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

// 整理静态数据放到monogoDB中
func main() {
	logging.Info("begin==")
	// FC数据库连接
	conn, err := dbr.Open("mysql", "finchina:finchina@tcp(114.55.105.11:3306)/finchina?charset=utf8", nil)
	if err != nil {
		logging.Debug("mysql onn", err)
	}
	sess := conn.NewSession(nil)
	// 个股
	StockTreatingData(sess)

	logging.Info("end==")
}

// 处理个股静态数据
func StockTreatingData(sess *dbr.Session) {
	logging.Info("stock begin==")
	// monogoDB 插入
	mgo_conn, err := mgo.DialWithTimeout(URL, time.Second*10)
	if err != nil {
		logging.Info("monogoDB 插入出错 %v", err)
	}
	mgo_collection := mgo_conn.DB("hgs").C(GLOBAL_SECRITY_TABLE)

	// 获取沪深股票信息
	secNm, err := new(fms.FcSecuNameTab).GetSecuNmList(sess)

	if err != nil {
		logging.Info("查询finance出错 %v", err)
	}
	// 遍历所有沪深股票
	for _, item := range secNm {
		var tss TagStockStatic
		// 根据证券id获取证券信息
		basinfo, err := new(stf.TQ_SK_BASICINFO).GetBasicinfoList(sess, item.SYMBOL.String)

		if err != nil {
			if err == dbr.ErrNotFound {
				logging.Info("查询证券信息未找到数据 %v", err)
			} else {
				logging.Info("查询证券信息err %v", err)
			}
		}
		// 转换证券代码
		swi := basinfo.EXCHANGE.String
		sym := basinfo.SYMBOL.String
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
		tss.SzSType = basinfo.SETYPE.String
		tss.SzStatus = basinfo.LISTSTATUS.String
		lde, err := strconv.Atoi(basinfo.LISTDATE.String)
		tss.NListDate = int32(lde)
		dse, err := strconv.Atoi(basinfo.DELISTDATE.String)
		tss.NDelistDate = int32(dse)
		// 根据公司内码获取股东信息
		shdn, errs := new(stf.TQ_SK_SHAREHOLDERNUM).GetSingleInfo(sess, item.COMPCODE.String)

		if errs != nil {
			if err == dbr.ErrNotFound {
				logging.Info("查询股东信息未找到数据 %v", err)
			} else {
				logging.Info("查询股东信息err %v", errs)
			}
		} else {
			cira, err := strconv.Atoi(shdn.CIRCSKAAMT.String)
			tss.LlCircuShare = int64(cira)
			tosa, err := strconv.Atoi(shdn.TOTALSHARE.String)
			tss.LlTotalShare = int64(tosa)
			tost, err := strconv.Atoi(shdn.TOTALSHAMT.String)
			tss.NHolders = int32(tost)
			if err != nil {
				logging.Info("类型转换 error %v", err)
			}
		}

		// 根据公司内码查询股票历史信息
		dklinfo, err := new(dkfm.Stock).GetSKTList5FC(sess, item.SECODE.String)

		if err != nil {
			if err == dbr.ErrNotFound {
				logging.Info("查询股票历史信息未找到数据 %v", err)
			} else {
				logging.Info("查询股票历史信息出错 error %v", err)
			}
		}
		var vol5 = 0
		for index, dkl := range dklinfo {
			if index == 0 {
				tss.NLastTradeDate = int32(dkl.TRADEDATE.Int64)
			}
			vol5 += int(dkl.VOL.Int64)
		}
		tss.LlLast5Volume = int64(vol5)
		// 查询公司业绩报表填充 每股收益和总资产
		tspe, err := new(stf.TQ_FIN_PROINCSTATEMENTNEW).GetSingleInfo(sess, item.COMPCODE.String)

		if err != nil {
			if err == dbr.ErrNotFound {
				logging.Info("查询公司业绩报表信息未找到数据 %v", err)
			} else {
				logging.Info("查询公司业绩报表信息出错 error %v", err)
			}
		}
		tss.NEPS = tspe.BASICEPS.Float64
		tss.LlTotalProfit = int64(tspe.TOTPROFIT.Float64)
		tss.LlNetProfit = int64(tspe.NETPROFIT.Float64)
		tss.NReportDate = int32(tspe.PUBLISHDATE.Int64)
		// 上市公司业绩快报 填充总资产
		tspce, err := new(stf.TQ_SK_PERFORMANCE).GetSingleInfo(sess, item.COMPCODE.String)

		if err != nil {
			if err == dbr.ErrNotFound {
				logging.Info("上市公司业绩快报信息未找到数据 %v", err)
			} else {
				logging.Info("查询上市公司业绩快报信息出错 error %v", err)
			}
		}
		tss.LlTotalProperty = tspce.TOTASSET.Float64
		//查询一般企业资产负债信息 填充流动资产
		tfp, err := new(stf.TQ_FIN_PROBALSHEETNEW).GetSingleInfo(sess, item.COMPCODE.String)

		if err != nil {
			if err == dbr.ErrNotFound {
				logging.Info("查询一般企业资产负债信息未找到数据 %v", err)
			} else {
				logging.Info("查询一般企业资产负债信息出错 error %v", err)
			}
		}
		tss.LlFlowProperty = int64(tfp.TOTCURRASSET.Float64)
		tss.NAVPS = int32(int64(tss.LlTotalProperty) / tss.LlTotalShare) // 总资产/总股本计算得到
		// 查询行业财务指标信息 填充主营业收入和利润
		tfi, err := new(stf.TQ_SK_BUSIINFO).GetSingleInfo(sess, item.COMPCODE.String)

		if err != nil {
			if err == dbr.ErrNotFound {
				logging.Info("查询行业财务指标信息未找到数据 %v", err)
			} else {
				logging.Info("查询行业财务指标信息出错 error %v", err)
			}
		}
		tss.LlMainIncoming = int64(tfi.TCOREBIZINCOME.Float64)
		tss.LlMainProfit = int64(tfi.TCOREBIZPROFIT.Float64)
		// 查询衍生财务指标信息 流动比率和速动比率
		tfpr, err := new(stf.TQ_FIN_PROINDICDATA).GetSingleInfo(sess, item.COMPCODE.String)

		if err != nil {
			if err == dbr.ErrNotFound {
				logging.Info("查询衍生财务指标信未找到数据 %v", err)
			} else {
				logging.Info("查询衍生财务指标信息出错 error %v", err)
			}
		}
		tss.NCurrentRatio = int32(tfpr.CURRENTRT.Float64)
		tss.NQuickMovingRatio = int32(tfpr.QUICKRT.Float64)

		_, err = mgo_collection.Upsert(bson.M{"nSID": tss.NSID}, &tss) // nSID存在时更新，不存在时插入，此语句效率较慢
		if err != nil {
			logging.Info("insert error %v", err)
		}
	}
	logging.Info("stock end==")
}
