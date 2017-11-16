package market

import (
	"time"

	"haina.com/market/hqpost/controllers/market/kline"
	"haina.com/market/hqpost/controllers/market/sidcode"

	. "haina.com/market/hqpost/controllers"
	"haina.com/market/hqpost/controllers/market/bline"
	"haina.com/market/hqpost/controllers/market/minline"

	"haina.com/market/hqpost/config"
	"haina.com/market/hqpost/models/redistore"
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

	//----------------------板块指数线----------------------/
	bline.NewBlockIndex().UpdateBlockIndexDayLine(cfg)
	//--------------------------------------------------/

	end := time.Now()
	logging.Info("Update Kline historical data successed, and running time:%v", end.Sub(start))

	redistore.NewHKLine(REDISKEY_HQPOST_EXECUTED_TIME).HQpostExecutedTime()
	/*********************结束时间***********************/
}
