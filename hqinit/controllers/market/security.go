package market

import (
	"time"

	"haina.com/market/hqinit/config"

	"haina.com/market/hqinit/controllers/market/security"
	"haina.com/share/logging"
)

func Update(cfg *config.AppConfig) {
	/*********************开始时间************************/
	start := time.Now()

	//股票代码表
	security.UpdateSecurityCodeTable()
	end := time.Now()
	logging.Info("Update Kline historical data successed, and running time:%v", end.Sub(start))

	//市场代码表及证券基本数据
	security.UpdateSecurityTable()

	//指数基本数据
	security.UpdateIndexTable()

	/*********************结束时间***********************/

}
