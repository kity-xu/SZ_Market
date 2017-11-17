package market

import (
	"os"
	"time"

	"haina.com/market/hqpost/config"
	"haina.com/market/hqpost/controllers/market/kline"
	"haina.com/market/hqpost/controllers/market/minline"
	"haina.com/market/hqpost/controllers/market/sidcode"
	"haina.com/share/logging"
)

func Update(cfg *config.AppConfig) {
	/*********************开始时间************************/
	start := time.Now()

	//股票代码表
	sids, err := sidcode.GetSecurityTable()
	if err != nil {
		return
	}

	//---------------------日周月年线---------------------/
	//K线
	if len(os.Args) == 3 {
		min := minline.NewMinKline(sids, cfg)
		kline.InitPath(cfg)

		switch os.Args[2] {
		case "-m":
			logging.Info("--------min----------")
			min.InitMinLine()
			// 分钟线
			min.HMinLine_1()
			min.HMinLine_5()
			min.HMinLine_15()
			min.HMinLine_30()
			min.HMinLine_60()
		case "-o":
			logging.Info("--------oDay----------")
			kline.HisDayKline(sids)
			kline.HisWeekKline(sids)
			kline.HisMonthKline(sids)
			kline.HisYearKline(sids)
		case "-day":
			logging.Info("--------day----------")
			kline.HisDayKline(sids)
		case "-week":
			logging.Info("--------week----------")
			kline.HisWeekKline(sids)
		case "-month":
			logging.Info("--------month----------")
			kline.HisMonthKline(sids)
		case "-year":
			logging.Info("--------year----------")
			kline.HisYearKline(sids)
		case "-m01":
			logging.Info("--------min01----------")
			min.InitMinLine()
			min.HMinLine_1()
		case "-m05":
			logging.Info("--------min05----------")
			min.InitMinLine()
			min.HMinLine_5()
		case "-m15":
			logging.Info("--------min15----------")
			min.InitMinLine()
			min.HMinLine_15()
		case "-m30":
			logging.Info("--------min30----------")
			min.InitMinLine()
			min.HMinLine_30()
		case "-m60":
			logging.Info("--------min60----------")
			min.InitMinLine()
			min.HMinLine_60()
		default:
			return
		}

	} else {
		logging.Info("--------all----------")
		kline.InitPath(cfg)
		kline.HisDayKline(sids)
		kline.HisWeekKline(sids)
		kline.HisMonthKline(sids)
		kline.HisYearKline(sids)
		//--------------------------------------------------/

		//----------------------分钟线-----------------------/
		min := minline.NewMinKline(sids, cfg)
		min.InitMinLine()
		// 分钟线
		min.HMinLine_1()
		min.HMinLine_5()
		min.HMinLine_15()
		min.HMinLine_30()
		min.HMinLine_60()
		//--------------------------------------------------/
	}
	//----------------------板块指数线----------------------/
	//bline.NewBlockIndex().UpdateBlockIndexDayLine(cfg)
	//--------------------------------------------------/

	end := time.Now()
	logging.Info("Update Kline historical data successed, and running time:%v", end.Sub(start))

	//redistore.NewHKLine(REDISKEY_HQPOST_EXECUTED_TIME).HQpostExecutedTime()
	/*********************结束时间***********************/
}
