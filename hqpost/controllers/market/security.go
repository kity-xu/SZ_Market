package market

import (
	"time"

	"haina.com/market/hqpost/controllers/market/sidcode"

	"haina.com/market/hqpost/controllers/market/kline"
	//"haina.com/market/hqpost/controllers/market/minline"

	"haina.com/market/hqpost/config"
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

	security := kline.NewSecurityKLine(sids, cfg)
	//--------------------------------------------------/
	//K线
	security.DayLine()
	security.WeekLine()
	security.MonthLine()
	security.YearLine()
	//--------------------------------------------------/

	//----------------------分钟线-----------------------/
	//	min := minline.NewMinKline(sids, cfg)
	//	// 分钟线
	//	min.HMinLine_1()
	//	min.HMinLine_5()
	//	min.HMinLine_15()
	//	min.HMinLine_30()
	//	min.HMinLine_60()

	end := time.Now()
	logging.Info("Update Kline historical data successed, and running time:%v", end.Sub(start))
	/*********************结束时间***********************/

}
